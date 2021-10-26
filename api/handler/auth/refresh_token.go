package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	jwtt "userland/api/jwt"
	"userland/store"

	"github.com/golang-jwt/jwt"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// refresh token handler still not working

func RefreshToken(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RefreshTokenRequest
		json.NewDecoder(r.Body).Decode(&request)

		refreshToken := request.RefreshToken

		token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
			}
			return []byte(os.Getenv("REFRESH_KEY")), nil
		})

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			res := map[string]string{
				"message": "refresh token expired",
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("went here")
			json.NewEncoder(w).Encode(ok)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			refreshUuid, ok := claims["refresh_uuid"].(string)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(ok)
				return
			}

			userId, _ := claims["user_id"].(string)
			deleted, delErr := tokenStore.DeleteUserId(refreshUuid)
			if delErr != nil || deleted == 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(delErr)
				return
			}

			ts, createErr := jwtt.GenerateToken(userId)
			if createErr != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(createErr)
				return
			}

			saveErr := jwtt.CreateAuth(userId, ts, tokenStore)
			if saveErr != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(createErr)
				return
			}

			cookie1 := &http.Cookie{
				Name:  "access_token",
				Value: ts.AccessToken,
			}

			cookie2 := &http.Cookie{
				Name:  "refresh_token",
				Value: ts.RefreshToken,
			}

			http.SetCookie(w, cookie1)
			http.SetCookie(w, cookie2)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		}

	}
}
