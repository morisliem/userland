package store

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetail struct {
	AccessUuid string
	UserId     string
	RefreshJti string
}

type RefreshDetail struct {
	AccessJti   string
	RefreshUuid string
	UserId      string
}

type TokenStore interface {
	StoreAccess(userId string, td TokenDetails) error
	StoreRefresh(userId string, td TokenDetails) error
	GetAtUserId(td *AccessDetail) (string, error)
	GetRtUserId(td *RefreshDetail) (string, error)
	SetEmailVerificationCode(email string, code int) error
	GetEmailVarificationCode(email string) (int, error)
	DeleteAtJti(uid string) (int64, error)
	DeleteRtJti(uid string) (int64, error)
}
