package rediss

import (
	"errors"
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
	// err := ts.db.Set(td.AccessUuid, userId, t.Sub(time.Now())).Err()
	// fmt.Println("access", td.AccessUuid)
	err := ts.db.Set(td.AccessUuid, userId, time.Until(time.Unix(td.AtExpires, 0))).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenStore) StoreRefresh(userId string, td store.TokenDetails) error {
	// fmt.Println("refresh", td.RefreshUuid)
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

func (ts *TokenStore) DeleteUserId(s string) (int64, error) {
	deleted, err := ts.db.Del(s).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (ts *TokenStore) SetEmailVerificationCode(email string, s int) error {

	err := ts.db.Set(email, s, time.Second*60).Err()
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

// func (ts *TokenStore) GetToken(ctx context.Context) error {
// 	return nil
// }
