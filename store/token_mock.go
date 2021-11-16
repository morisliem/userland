package store

import "github.com/stretchr/testify/mock"

type MockTokenStore struct {
	mock.Mock
}

func (m *MockTokenStore) StoreAccess(uid string, td TokenDetails) error {
	args := m.Called(uid, td)
	return args.Error(0)
}

func (m *MockTokenStore) StoreRefresh(uid string, td TokenDetails) error {
	args := m.Called(uid, td)
	return args.Error(0)
}

func (m *MockTokenStore) GetAtUserId(td *AccessDetail) (string, error) {
	args := m.Called(td)
	return args.String(0), args.Error(1)
}

func (m *MockTokenStore) GetRtUserId(td *RefreshDetail) (string, error) {
	args := m.Called(td)
	return args.String(0), args.Error(1)
}

func (m *MockTokenStore) SetEmailVerificationCode(uid string, code int) error {
	args := m.Called(uid, code)
	return args.Error(0)
}

func (m *MockTokenStore) GetEmailVarificationCode(uid string) (int, error) {
	args := m.Called(uid)
	return args.Int(0), args.Error(1)
}

func (m *MockTokenStore) SetNewEmail(uid string, email string) error {
	args := m.Called(uid, email)
	return args.Error(0)
}

func (m *MockTokenStore) GetNewEmail(uid string) (string, error) {
	args := m.Called(uid)
	return args.String(0), args.Error(1)
}

func (m *MockTokenStore) DeleteAtJti(uid string) (int64, error) {
	args := m.Called(uid)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockTokenStore) DeleteRtJti(uid string) (int64, error) {
	args := m.Called(uid)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockTokenStore) HasRefreshToken(jti string) (bool, error) {
	args := m.Called(jti)
	return args.Bool(0), args.Error(1)
}
