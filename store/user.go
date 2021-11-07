package store

import (
	"context"
	"time"
)

type User struct {
	Id         string
	Fullname   string
	Email      string
	Password   string
	Location   string
	Bio        string
	Web        string
	Picture    string
	VerCode    int
	Created_at time.Time
}

type UserInfo struct {
	UserId    string
	SessionId string
	Name      string
}

type UserSession struct {
	Is_current bool
	Ip         string
	Client     []UserInfo
	Created_at time.Time
	Updated_at time.Time
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
	EmailActive(ctx context.Context, u User) (bool, error)
	GetUserDetail(ctx context.Context, uid string) (User, error)
	SetUserPicture(ctx context.Context, uid string, pict string) error
	DeleteUserPicture(ctx context.Context, uid string) error
	UpdateUserDetail(ctx context.Context, u User, uid string) error
	GetUserEmail(ctx context.Context, uid string) (User, error)
	DeleteAccount(ctx context.Context, uid string) error
	EmailExist(ctx context.Context, u User) error
	SetUserSession(ctx context.Context, t TokenDetails, uid string, ip string, device string) error
	GetUserSession(ctx context.Context, uid string, sessionId string) (UserSession, error)
	UpdateUserSession(ctx context.Context, sessionId string) error
	DeleteCurrentSession(ctx context.Context, sessionId string) error
	DeleteOtherSession(ctx context.Context, uid string, sessionId string) error
	GetSessionsId(ctx context.Context, uid string) ([]string, error)
	GetUserProfilePicture(ctx context.Context, uid string) (string, error)
}
