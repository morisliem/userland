package jwt

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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

func ExtractAccessTokenMetadata(r *http.Request) (string, error) {
	tkn := ExtractToken(r)

	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_KEY")), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return "", err
		}
		return accessUuid, nil
	}
	return "", nil
}

func ExtractRefreshTokenMetadata(r *http.Request) (string, error) {
	tkn := ExtractToken(r)

	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_KEY")), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		access_jti := claims["access_jti"].(string)
		return access_jti, nil
	}
	return "", nil
}
