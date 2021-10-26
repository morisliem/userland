package store

import (
	"context"
	"time"
)

type User struct {
	Id         uint64    `json:"Id"`
	Fullname   string    `json:"fullname"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Location   string    `json:"location"`
	Bio        string    `json:"bio"`
	Web        string    `json:"web"`
	Picture    string    `json:"picture"`
	VerCode    int       `json:"verification_code"`
	Created_at time.Time `json:"created_at"`
}

type UserStore interface {
	GetUserId(ctx context.Context, u User) (string, error)
	RegisterUser(ctx context.Context, u User) error
	ResetPassword(ctx context.Context, uid string, u User) error
	ValidateCode(ctx context.Context, u User) error
}
