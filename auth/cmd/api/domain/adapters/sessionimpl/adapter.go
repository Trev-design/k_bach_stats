package sessionimpl

import "github.com/google/uuid"

type Adapter interface {
	GetVerifyData(cookie string) (uuid.UUID, string, string, error)
	DeleteVerifySession(id string) error
	SetVerifyData(accountID string) (string, string, error)
	SetRefreshData(accountID, ip, userAgent string) (string, error)
	VerifyRefreshData(cookie, ip, userAgent string) (string, error)
	RemoveRefreshData(cookie, ip, userAgent string) error

	CloseSession() error
	HandleBackground()
}
