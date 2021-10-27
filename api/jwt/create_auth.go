package jwt

import (
	"userland/store"
)

func CreateATAuth(userId string, td store.TokenDetails, ts store.TokenStore) error {

	errAccess := ts.StoreAccess(userId, td)
	if errAccess != nil {
		return errAccess
	}

	return nil
}

func CreateRTAuth(userId string, td store.TokenDetails, ts store.TokenStore) error {

	errRefresh := ts.StoreRefresh(userId, td)
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func FetchATAuth(auth *store.AccessDetail, ts store.TokenStore) (string, error) {
	userId, err := ts.GetAtUserId(auth)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func FetchRTAuth(auth *store.RefreshDetail, ts store.TokenStore) (string, error) {
	userId, err := ts.GetRtUserId(auth)
	if err != nil {
		return "", err
	}
	return userId, nil
}
