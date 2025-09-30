package dbimpl

import (
	"auth_server/cmd/api/domain/types"

	"github.com/google/uuid"
)

type Adapter interface {
	// cleans up the database instance based on certain scenarions like server shutdown or database reconnect
	CloseDatabase() error

	// adds a new accout default is unauthorized and has no abo plan till the user is verified
	AddUser(newAccount *types.NewAccountDM) (string, error)

	// gets a account based on its id
	GetUser(id uuid.UUID) (*types.AccountDM, error)

	// gets an account based on its email. to be able to signs the account in for example
	GetUserByEmail(email string) (*types.AccountDM, error)

	// verifies the account and changes its state to verified and his abo_type which is COMMUNITY
	UpdateState(id uuid.UUID) (*types.AccountDM, error)

	// to change the password if you forgot your password or just want to change it
	ChangePassword(id uuid.UUID, passwordHash string) (*types.AccountDM, error)
}
