package dbimpl

import (
	"auth_server/cmd/api/domain/types"

	"github.com/google/uuid"
)

type Adapter interface {
	CloseDatabase() error
	AddUser(newAccount *types.NewAccountDM) (string, error)
	GetUser(id uuid.UUID) (*types.AccountDM, error)
	GetUserByEmail(email string) (*types.AccountDM, error)
	UpdateState(id uuid.UUID) (*types.AccountDM, error)
	ChangePassword(id uuid.UUID, passwordHash string) (*types.AccountDM, error)
}
