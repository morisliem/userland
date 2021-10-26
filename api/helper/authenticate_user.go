package helper

import (
	"net/http"
	"userland/api/jwt"
	"userland/store"
)

func AuthenticateUser(r *http.Request, tokenStore store.TokenStore) (string, error) {
	tokenAuth, err := jwt.ExtractTokenMetadata(r)
	if err != nil {
		return "", err
	}

	userId, err := jwt.FetchAuth(tokenAuth, tokenStore)
	if err != nil {
		return "", err
	}

	return userId, nil
}
