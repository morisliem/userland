package store

import (
	"context"
	"time"
)

type User struct {
	Id         string    `json:"Id"`
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

type UserInfo struct {
	SessionId string `json:"id"`
	Name      string `json:"name"`
}

type UserSession struct {
	Is_current bool       `json:"is_current"`
	Ip         string     `json:"ip"`
	Client     []UserInfo `json:"clients"`
	Created_at time.Time  `json:"created_at"`
	Updated_at time.Time  `json:"updated_at"`
}

type UserStore interface {
	GetUserId(ctx context.Context, u User) (string, error)
	GetUserid(ctx context.Context, email string) (string, error)
	RegisterUser(ctx context.Context, u User, rn int) error
	UpdatePassword(ctx context.Context, uid string, u User) error
	ValidateCode(ctx context.Context, u User) error
	GetUserCode(ctx context.Context, u User) (int, error)
	GetPassword(ctx context.Context, uid string) (string, error)
	GetPasswords(ctx context.Context, uid string) ([]string, error)
	GetUserState(ctx context.Context, u User) (int, error)
	GetUserDetail(ctx context.Context, uid string) (User, error)
	UpdateUserDetail(ctx context.Context, u User, uid string) error
	GetUserEmail(ctx context.Context, uid string) (User, error)
	UpdateUserEmail(ctx context.Context, u User, uid string) error
	DeleteAccount(ctx context.Context, uid string) error
	EmailExist(ctx context.Context, u User) error
	SetUserSession(ctx context.Context, t TokenDetails, uid string, ip string) error
	GetUserSession(ctx context.Context, uid string) (UserSession, error)
}
