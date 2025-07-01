package domainimpl

import "auth_server/cmd/api/domain/types"

type Adapter interface {
	Register(account *types.NewAccountRequestDTO) (*types.VerifySessionDTO, error)
	Verify(verifyRequest *types.VerifyAccountDTO) (*types.NewAccountSessionDTO, error)
	Login(loginRequest *types.LoginAccountDTO) (*types.NewAccountSessionDTO, error)
	NewPassword(email string) (*types.VerifySessionDTO, error)
	ChangePassword(changePassRequest *types.ChangePasswordDTO) (*types.NewAccountSessionDTO, error)
	RefreshSession(creds *types.RefreshSessionDTO) (*types.NewAccountSessionDTO, error)
	RemoveSession(creds *types.RemoveSessionDTO) error
}

type Instance interface {
	ComputeBackgroundServices()
	ListenForShutdown()
}
