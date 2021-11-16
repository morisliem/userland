package store

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type TokenStore interface {
	StoreAccess(userId string, td TokenDetails) error
	StoreRefresh(userId string, td TokenDetails) error
	GetAtUserId(atJti string) (string, error)
	GetRtUserId(atJti string) (string, error)
	SetEmailVerificationCode(uid string, code int) error
	GetEmailVarificationCode(uid string) (int, error)
	SetNewEmail(uid string, email string) error
	GetNewEmail(uid string) (string, error)
	DeleteAtJti(atJti string) (int64, error)
	DeleteRtJti(atJti string) (int64, error)
	HasRefreshToken(jti string) (bool, error)
}
