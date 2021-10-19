package postgres

import (
	"database/sql"
	"context"
	"userland/store"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) store.UserStore {
	return &UserStore {
		db: db,
	}
}

func (us *UserStore) GetUser(ctx context.Context) error{
	_, _ = us.db.QueryContext(ctx, "")
	return nil
}