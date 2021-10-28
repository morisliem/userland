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

func (ts *TokenStore) StoreAccess(userId string, td store.TokenDetails) error {
	err := ts.db.Set(td.AccessUuid, userId, time.Until(time.Unix(td.AtExpires, 0))).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenStore) StoreRefresh(userId string, td store.TokenDetails) error {
	err := ts.db.Set(td.RefreshUuid, userId, time.Until(time.Unix(td.RtExpires, 0))).Err()
	if err != nil {
		return err
	}
	return nil
}

func (ts *TokenStore) GetAtUserId(td *store.AccessDetail) (string, error) {
	userId, err := ts.db.Get(td.AccessUuid).Result()
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (ts *TokenStore) GetRtUserId(td *store.RefreshDetail) (string, error) {
	userId, err := ts.db.Get(td.RefreshUuid).Result()
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (ts *TokenStore) DeleteJti(s string) (int64, error) {
	deleted, err := ts.db.Del(s).Result()
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

	code, err := strconv.Atoi(res)
	if err != nil {
		return 0, err
	}
	return code, nil
}
