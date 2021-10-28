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
	DeleteJti(jwtId string) (int64, error)
	SetEmailVerificationCode(email string, code int) error
	GetEmailVarificationCode(email string) (int, error)
	// GetToken(ctx context.Context) error
}
