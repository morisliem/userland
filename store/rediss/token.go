package rediss

import (
	"errors"
	"os"
	"strconv"
	"time"
	"userland/store"

	"github.com/go-redis/redis"
)

type TokenStore struct {
	db *redis.Client
}

func NewTokenStore(db *redis.Client) store.TokenStore {
	return &TokenStore{
		db: db,
	}
}

// Idea so far
/*

Creating refresh token jti when creating the access token
So the next time user ask for refresh token, that refresh token jti will be its id

Than store the access token in redis as "access_token:access uid" is the key
Then store the refresh token in redis as "refresh_token:access uid" is the key

That way, i can remove the other session without knowing the refresh token jti

The loophole is when user has not request for refresh token, the access token is already have its refresh token jti
This issue can be solved by creating another function to check if the access token has been updated or not
If it's been updated, meaning that the refresh token has also been created
Otherwise, don't need to remove the refresh token since it's not been created just yet


*/

func (ts *TokenStore) StoreAccess(userId string, td store.TokenDetails) error {
	key := "access_token:" + td.AccessUuid

	err := ts.db.Set(key, userId, time.Until(time.Unix(td.AtExpires, 0))).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenStore) StoreRefresh(userId string, td store.TokenDetails) error {
	key := "refresh_token:" + td.AccessUuid

	err := ts.db.Set(key, userId, time.Until(time.Unix(td.RtExpires, 0))).Err()
	if err != nil {
		return err
	}
	return nil
}

func (ts *TokenStore) GetAtUserId(td *store.AccessDetail) (string, error) {
	key := "access_token:" + td.AccessUuid

	res, err := ts.db.Get(key).Result()
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("token is expired")
	}
	return td.UserId, nil
}

func (ts *TokenStore) GetRtUserId(td *store.RefreshDetail) (string, error) {
	key := "refresh_token:" + td.AccessJti

	res, err := ts.db.Get(key).Result()
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("token is expired")
	}
	return td.UserId, nil
}

func (ts *TokenStore) DeleteAtJti(atJti string) (int64, error) {
	key := "access_token:" + atJti
	deleted, err := ts.db.Del(key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (ts *TokenStore) DeleteRtJti(atJti string) (int64, error) {
	key := "refresh_token:" + atJti
	deleted, err := ts.db.Del(key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (ts *TokenStore) SetEmailVerificationCode(email string, s int) error {
	duration, _ := strconv.Atoi(os.Getenv("EMAIL_CODE_DURATION"))
	err := ts.db.Set(email, s, time.Second*time.Duration(duration)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenStore) GetEmailVarificationCode(email string) (int, error) {
	res, err := ts.db.Get(email).Result()
	if err != nil {
		return 0, errors.New("code is expired")
	}
	if len(res) == 0 {
		return 0, errors.New("token is expired")
	}

	code, err := strconv.Atoi(res)
	if err != nil {
		return 0, err
	}
	return code, nil
}
