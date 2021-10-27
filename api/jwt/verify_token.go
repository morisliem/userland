package jwt

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func VerifyAccessToken(r *http.Request) (*jwt.Token, error) {
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
	return token, nil
}

func VerifyRefreshToken(r *http.Request) (*jwt.Token, error) {
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
	return token, nil
}
