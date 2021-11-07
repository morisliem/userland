package jwt

import "userland/store"

// Deleting the token in redis
func DeleteATAuth(atJti string, ts store.TokenStore) (int64, error) {
	del, err := ts.DeleteAtJti(atJti)
	if err != nil {
		return 0, err
	}

	return del, nil
}

// Deleting the token in redis
func DeleteRTAuth(atJti string, ts store.TokenStore) (int64, error) {
	del, err := ts.DeleteRtJti(atJti)
	if err != nil {
		return 0, err
	}

	return del, nil
}
