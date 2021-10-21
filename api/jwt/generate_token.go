package jwt

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userId string) (string, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["autorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString([]byte(os.Getenv("ACCESS_KEY")))
	if err != nil {
		return "", err
	}
	return token, nil
}
