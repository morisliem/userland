package jwt

import (
	"userland/store"
)

func CreateAuth(userId string, td store.TokenDetails, ts store.TokenStore) error {

	errAccess := ts.StoreAccess(userId, td)
	if errAccess != nil {
		return errAccess
	}

	errRefresh := ts.StoreRefresh(userId, td)
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func FetchAuth(auth *store.AccessDetail, ts store.TokenStore) (string, error) {
	userId, err := ts.GetUserId(auth)
	if err != nil {
		return "", err
	}
	return userId, nil
}
