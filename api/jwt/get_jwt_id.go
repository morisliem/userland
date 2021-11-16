package jwt

import (
	"net/http"
)

func GetAtJtiFromAccessToken(r *http.Request) (string, error) {
	atJti, err := ExtractAccessTokenMetadata(r)
	if err != nil {
		return "", err
	}

	return atJti, nil
}

func GetAtJtiFromRefreshToken(r *http.Request) (string, error) {
	atJti, err := ExtractRefreshTokenMetadata(r)
	if err != nil {
		return "", err
	}

	return atJti, nil
}
