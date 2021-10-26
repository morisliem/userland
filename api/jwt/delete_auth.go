package jwt

import "userland/store"

func DeleteAuth(userId string, ts store.TokenStore) (int64, error) {
	del, err := ts.DeleteUserId(userId)
	if err != nil {
		return 0, err
	}

	return del, nil
}
