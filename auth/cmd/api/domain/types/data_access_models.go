package types

type RefreshSessionDAO struct {
	AccountID string
	IPAddress string
	UserAgent string
}

type JWTDAO struct {
	ID   string
	Role string
}
