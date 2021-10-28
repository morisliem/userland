package jwt

import "userland/store"

func DeleteATAuth(jwtId string, ts store.TokenStore) (int64, error) {
	del, err := ts.DeleteJti(jwtId)
	if err != nil {
		return 0, err
	}

	return del, nil
}

func DeleteRTAuth(jwtId string, ts store.TokenStore) (int64, error) {
	del, err := ts.DeleteJti(jwtId)
	if err != nil {
		return 0, err
	}

	return del, nil
}
