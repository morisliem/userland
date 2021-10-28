package jwt

import (
	"net/http"
)

func GetAtJtinRtJti(r *http.Request) (string, string, error) {
	tokenAuth, err := ExtractAccessTokenMetadata(r)
	if err != nil {
		return "", "", err
	}

	return tokenAuth.AccessUuid, tokenAuth.RefreshJti, nil
}

func GetAtJtiFromRt(r *http.Request) (string, error) {
	tokenAuth, err := ExtractRefreshTokenMetadata(r)
	if err != nil {
		return "", err
	}
	return tokenAuth.AccessJti, nil
}

func GetRtJti(r *http.Request) (string, error) {
	tokenAuth, err := ExtractRefreshTokenMetadata(r)
	if err != nil {
		return "", err
	}

	return tokenAuth.RefreshUuid, nil
}
