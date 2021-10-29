package helper

import (
	"net/http"
	"userland/api/jwt"
	"userland/store"
)

func AuthenticateUserAccessToken(r *http.Request, tokenStore store.TokenStore) (string, error) {
	tokenAuth, err := jwt.ExtractAccessTokenMetadata(r)
	if err != nil {
		return "", err
	}

	//IsAtStillActive
	userId, err := jwt.FetchATAuth(tokenAuth, tokenStore)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func AuthenticateUserRefreshToken(r *http.Request, tokenStore store.TokenStore) (string, error) {
	tokenAuth, err := jwt.ExtractRefreshTokenMetadata(r)
	if err != nil {
		return "", err
	}

	//IsRtStillActive
	userId, err := jwt.FetchRTAuth(tokenAuth, tokenStore)
	if err != nil {
		return "", err
	}

	return userId, nil
}
