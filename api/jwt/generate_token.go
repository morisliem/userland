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

func GenerateAccessToken(userId string, atJti string, rtJti string, ts store.TokenStore) (store.TokenDetails, error) {
	td := store.TokenDetails{}
	atDuration, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION"))

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = userId

	if atJti == "" {
		td.AtExpires = time.Now().Add(time.Minute * time.Duration(atDuration)).Unix()
		accessUuid, _ := uuid.NewV4()
		td.AccessUuid = fmt.Sprintf("%v", accessUuid)
		atClaims["access_uuid"] = td.AccessUuid

	} else {
		td.AtExpires = time.Now().Add(time.Minute * time.Duration(atDuration)).Unix()
		td.AccessUuid = atJti
		atClaims["access_uuid"] = atJti
	}

	if rtJti == "" {
		refresh_jti, _ := uuid.NewV4()
		td.RefreshUuid = fmt.Sprintf("%v", refresh_jti)
		atClaims["refresh_jti"] = td.RefreshUuid

	} else {
		atClaims["refresh_jti"] = rtJti
	}

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

func GenerateRefreshToken(userId string, atJti string, rtJti string, ts store.TokenStore) (store.TokenDetails, error) {
	td := store.TokenDetails{}
	rtDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION"))
	td.RtExpires = time.Now().Add(time.Minute * time.Duration(rtDuration)).Unix()
	td.RefreshUuid = rtJti
	td.AccessUuid = atJti

	var err error

	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userId
	rtClaims["refresh_uuid"] = td.RefreshUuid
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
