package jwt

import (
	"fmt"
	"os"
	"time"
	"userland/store"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(userId string) (store.TokenDetails, error) {
	td := store.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 30).Unix()
	accessUuid, _ := uuid.NewV4()
	td.AccessUuid = fmt.Sprintf("%v", accessUuid)

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = userId
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_KEY")))
	if err != nil {
		return td, err
	}
	return td, nil
}

func GenerateRefreshToken(userId string) (store.TokenDetails, error) {
	td := store.TokenDetails{}
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	refreshUuid, _ := uuid.NewV4()
	td.RefreshUuid = fmt.Sprintf("%v", refreshUuid)
	var err error
	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userId
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_KEY")))
	if err != nil {
		return td, err
	}
	return td, nil
}
