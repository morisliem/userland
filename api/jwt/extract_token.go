package jwt

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"userland/store"

	"github.com/golang-jwt/jwt"
)

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	splitBearerToken := strings.Split(bearerToken, " ")

	if len(splitBearerToken) == 2 {
		return splitBearerToken[1]
	}

	return ""
}

func ExtractAccessTokenMetadata(r *http.Request) (*store.AccessDetail, error) {
	tkn := ExtractToken(r)

	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId := claims["user_id"].(string)
		refresh_jti := claims["refresh_jti"].(string)
		return &store.AccessDetail{
			AccessUuid: accessUuid,
			UserId:     userId,
			RefreshJti: refresh_jti,
		}, nil
	}
	return nil, nil
}

func ExtractRefreshTokenMetadata(r *http.Request) (*store.RefreshDetail, error) {
	tkn := ExtractToken(r)

	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId := claims["user_id"].(string)
		access_jti := claims["access_jti"].(string)
		return &store.RefreshDetail{
			RefreshUuid: refreshUuid,
			UserId:      userId,
			AccessJti:   access_jti,
		}, nil
	}
	return nil, nil
}
