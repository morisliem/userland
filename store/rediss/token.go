package rediss

import (
	"errors"
	"fmt"
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
		fmt.Println("if error occur")
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

func (ts *TokenStore) GetAtUserId(atJti string) (string, error) {
	key := "access_token:" + atJti

	res, err := ts.db.Get(key).Result()
	if len(res) == 0 {
		return "", errors.New("token is expired")
	}

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return "", err
	}

	return res, nil
}

func (ts *TokenStore) GetRtUserId(atJti string) (string, error) {
	key := "refresh_token:" + atJti

	res, err := ts.db.Get(key).Result()
	if len(res) == 0 {
		return "", errors.New("token is expired")
	}

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return "", err
	}
	return res, nil
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

	return res, nil
}
