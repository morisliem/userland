package mysession

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"userland/api/helper"
	"userland/api/jwt"
	"userland/api/response"
	"userland/store"
)

type GetATResponse struct {
	Value      string    `json:"value"`
	Type       string    `json:"type"`
	Expired_at time.Time `json:"expired_at"`
}

func GetAccessToken(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helper.AuthenticateUserRefreshToken(r, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		atJti, err := jwt.GetAtJtiFromRt(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		rtJti, err := jwt.GetRtJti(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		// generated new access token by adding the current jwt id (access token) and adding refresh token jwt id
		res, err := jwt.GenerateAccessToken(userId, atJti, rtJti)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		saveErr := jwt.CreateATAuth(userId, res, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(saveErr.Error()))
			return
		}

		err = userStore.UpdateUserSession(r.Context(), atJti)
		if err != nil {
			fmt.Println("here ;(")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newToken := &GetATResponse{
			Value:      res.AccessToken,
			Expired_at: time.Unix(res.AtExpires, 0),
			Type:       "jwt",
		}

		response := map[string]GetATResponse{
			"access_token": *newToken,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
