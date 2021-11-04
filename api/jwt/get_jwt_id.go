package jwt

import (
	"net/http"
)

func GetAtJtiNRtJtiFromAccessToken(r *http.Request) (string, string, error) {
	tokenAuth, err := ExtractAccessTokenMetadata(r)
	if err != nil {
		return "", "", err
	}

	return tokenAuth.AccessUuid, tokenAuth.RefreshJti, nil
}

func GetAtJtiNRtJtiFromRefreshToken(r *http.Request) (string, string, error) {
	tokenAuth, err := ExtractRefreshTokenMetadata(r)
	if err != nil {
		return "", "", err
	}

	return tokenAuth.AccessJti, tokenAuth.RefreshUuid, nil
}
