package store

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) GetUserId(ctx context.Context, u User) (string, error) {
	args := m.Called(ctx, u)
	return args.String(0), args.Error(1)
}

func (m *MockUserStore) GetUserid(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *MockUserStore) RegisterUser(ctx context.Context, u User, rn int) error {
	args := m.Called(ctx, u, rn)
	return args.Error(0)
}

func (m *MockUserStore) UpdatePassword(ctx context.Context, uid string, u User) error {
	args := m.Called(ctx, uid, u)
	return args.Error(0)
}

func (m *MockUserStore) ValidateCode(ctx context.Context, u User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserStore) GetUserCode(ctx context.Context, u User) (int, error) {
	args := m.Called(ctx, u)
	return args.Int(0), args.Error(1)
}

func (m *MockUserStore) GetPassword(ctx context.Context, uid string) (string, error) {
	args := m.Called(ctx, uid)
	return args.String(0), args.Error(1)
}

func (m *MockUserStore) GetPasswords(ctx context.Context, uid string) ([]string, error) {
	args := m.Called(ctx, uid)

	var list []string
	list = append(list, args.String(0))

	return list, args.Error(1)
}

func (m *MockUserStore) EmailActive(ctx context.Context, u User) (bool, error) {
	args := m.Called(ctx, u)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserStore) GetUserDetail(ctx context.Context, uid string) (User, error) {
	args := m.Called(ctx, uid)
	var newUser User

	return newUser, args.Error(1)
}

func (m *MockUserStore) SetUserPicture(ctx context.Context, uid string, pict string) error {
	args := m.Called(ctx, uid, pict)
	return args.Error(0)
}

func (m *MockUserStore) DeleteUserPicture(ctx context.Context, uid string) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockUserStore) UpdateUserDetail(ctx context.Context, u User, uid string) error {
	args := m.Called(ctx, u, uid)
	return args.Error(0)
}

func (m *MockUserStore) GetUserEmail(ctx context.Context, uid string) (User, error) {
	args := m.Called(ctx, uid)
	var newUser User
	return newUser, args.Error(1)
}

func (m *MockUserStore) DeleteAccount(ctx context.Context, uid string) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockUserStore) EmailExist(ctx context.Context, u User) error {
	args := m.Called(ctx, u)

	return args.Error(0)
}

func (m *MockUserStore) SetUserSession(ctx context.Context, t TokenDetails, uid string, ip string, device string) error {
	args := m.Called(ctx, t, uid, ip, device)
	return args.Error(0)
}

func (m *MockUserStore) GetUserSession(ctx context.Context, uid string, sessionId string) (UserSession, error) {
	args := m.Called(ctx, uid, sessionId)
	var newSession UserSession
	return newSession, args.Error(0)
}

func (m *MockUserStore) UpdateUserSession(ctx context.Context, sessionId string) error {
	args := m.Called(ctx, sessionId)
	return args.Error(0)
}

func (m *MockUserStore) DeleteCurrentSession(ctx context.Context, sessionId string) error {
	args := m.Called(ctx, sessionId)
	return args.Error(0)
}

func (m *MockUserStore) DeleteOtherSession(ctx context.Context, uid string, sessionId string) error {
	args := m.Called(ctx, uid, sessionId)
	return args.Error(0)
}

func (m *MockUserStore) GetSessionsId(ctx context.Context, uid string) ([]string, error) {
	args := m.Called(ctx, uid)

	var list []string
	list = append(list, args.String(0))

	return list, args.Error(1)
}

func (m *MockUserStore) GetUserProfilePicture(ctx context.Context, uid string) (string, error) {
	args := m.Called(ctx, uid)
	return args.String(0), args.Error(1)
}
