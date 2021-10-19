package store

import (
	"context"
)

type UserStore interface {
	GetUser(ctx context.Context) error
}