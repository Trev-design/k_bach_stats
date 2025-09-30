package domainimpl

import "auth_server/cmd/api/domain/types"

// the bridge between api and business logic
type Adapter interface {
	// tries to register a new account. if the registration failed you get an error
	Register(account *types.NewAccountRequestDTO) (*types.VerifySessionDTO, error)

	// tries to verify a registered account. if the registration failed you get an error
	Verify(verifyRequest *types.VerifyAccountDTO) (*types.NewAccountSessionDTO, error)

	// tries to login an verified and registered account. if the login failed you get no session but an error
	Login(loginRequest *types.LoginAccountDTO) (*types.NewAccountSessionDTO, error)

	// tries to register a password change. if the registration failed you get an error
	NewPassword(email string) (*types.VerifySessionDTO, error)

	// tries to change a password of an account which is authorized to change its password. if the change failed you get an error
	ChangePassword(changePassRequest *types.ChangePasswordDTO) (*types.NewAccountSessionDTO, error)

	// tries to refresh a session of a users account which has access to a session. if the refresh failed you get an error
	RefreshSession(creds *types.RefreshSessionDTO) (*types.NewAccountSessionDTO, error)

	// tries to remove a session of a users account which has access to a session because he logged out or the session is expired for example.
	// if the remove failed you get an error.
	RemoveSession(creds *types.RefreshSessionDTO) error
}

// the inteface to control the lifetime of the business logic
type Instance interface {
	// registers and compute all background processes of your business logic
	ComputeBackgroundServices()

	// cleans up all services of the business logic when the server shuts down
	ListenForShutdown()
}
