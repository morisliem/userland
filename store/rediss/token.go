package rediss

import (
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

// func (ts *TokenStore) GetToken(ctx context.Context) error {
// 	return nil
// }
