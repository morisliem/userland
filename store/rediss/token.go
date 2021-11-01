package rediss

import (
	"errors"
	"os"
	"strconv"
	"time"
	"userland/store"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
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

func (ts *TokenStore) HasRefreshToken(jti string) (bool, error) {
	key := "refresh_token:" + jti

	res, _ := ts.db.Get(key).Result()

	if len(res) == 0 {
		return false, nil
	} else {
		return true, nil
	}
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

func (ts *TokenStore) SetEmailVerificationCode(uid string, s int) error {
	duration, _ := strconv.Atoi(os.Getenv("EMAIL_CODE_DURATION"))
	key := "code:" + uid
	err := ts.db.Set(key, s, time.Second*time.Duration(duration)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenStore) GetEmailVarificationCode(uid string) (int, error) {
	key := "code:" + uid
	res, err := ts.db.Get(key).Result()
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return 0, err
	}
	if len(res) == 0 {
		return 0, errors.New("code is expired")
	}

	code, err := strconv.Atoi(res)
	if err != nil {
		return 0, err
	}
	return code, nil
}

func (ts *TokenStore) SetNewEmail(uid string, email string) error {
	key := "email:" + uid
	duration, _ := strconv.Atoi(os.Getenv("EMAIL_CODE_DURATION"))
	err := ts.db.Set(key, email, time.Second*time.Duration(duration)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenStore) GetNewEmail(uid string) (string, error) {
	key := "email:" + uid
	res, err := ts.db.Get(key).Result()
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("code is expired")
	}

	return res, nil
}
