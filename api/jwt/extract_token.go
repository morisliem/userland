package jwt

import (
	"net/http"
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

func ExtractTokenMetadata(r *http.Request) (*store.AccessDetail, error) {
	token, err := VerifyToken(r)
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
		return &store.AccessDetail{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, nil
}
