package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tkn := ExtractToken(r)

	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN")), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) (string, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return "", err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return "", errors.New("token is invalid")
	}

	tmp := fmt.Sprintf("%v", token.Claims)
	loc := strings.Index(tmp, "user_id:")
	userId := tmp[loc+8 : len(tmp)-1]

	return userId, nil
}
