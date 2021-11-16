package jwt

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"userland/store"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(userId string, ts store.TokenStore) (store.TokenDetails, error) {
	td := store.TokenDetails{}
	atDuration, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION"))
	td.AtExpires = time.Now().Add(time.Minute * time.Duration(atDuration)).Unix()
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

	errAccess := ts.StoreAccess(userId, td)
	if errAccess != nil {
		return td, errAccess
	}

	return td, nil
}

func GenerateRefreshToken(userId string, atJti string, ts store.TokenStore) (store.TokenDetails, error) {
	td := store.TokenDetails{}
	rtDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION"))
	td.RtExpires = time.Now().Add(time.Minute * time.Duration(rtDuration)).Unix()
	td.AccessUuid = atJti

	var err error

	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userId
	rtClaims["access_jti"] = td.AccessUuid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_KEY")))
	if err != nil {
		return td, err
	}

	errAccess := ts.StoreRefresh(userId, td)
	if errAccess != nil {
		return td, errAccess
	}
	return td, nil
}
