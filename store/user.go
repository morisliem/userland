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
	Created_at time.Time `json:"created_at"`
}

type UserStore interface {
	GetUser(ctx context.Context) error
	RegisterUser(ctx context.Context, u User) error
}
