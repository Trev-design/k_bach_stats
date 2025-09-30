package impl

import (
	"auth_server/cmd/api/domain/types"
	"errors"
)

// tries to register a new account
func (impl *Impl) Register(account *types.NewAccountRequestDTO) (*types.VerifySessionDTO, error) {
	if !passwordConfirmed(account.Password, account.Confirmation) {
		return nil, errors.New("invalid credentials")
	}

	passwordHash, err := getHashedPassword(account.Password)
	if err != nil {
		return nil, err
	}

	if !isValidEmail(account.Email) {
		return nil, errors.New("invalid credentials")
	}

	id, err := impl.db.AddUser(&types.NewAccountDM{
		Name:         account.Name,
		Email:        account.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, err
	}

	number, cookie, err := impl.session.SetVerifyData(id)
	if err != nil {
		return nil, err
	}

	if err := impl.sendVerifyData(number, account.Email, account.Name); err != nil {
		return nil, err
	}

	return &types.VerifySessionDTO{Name: account.Name, Cookie: cookie}, nil
}

// tries to verify an account
func (impl *Impl) Verify(verifyRequest *types.VerifyAccountDTO) (*types.NewAccountSessionDTO, error) {
	id, sessionID, number, err := impl.session.GetVerifyData(verifyRequest.Cookie)
	if err != nil {
		return nil, err
	}

	if number != verifyRequest.VerifyCode {
		return nil, errors.New("invalid verify code")
	}

	if err := impl.session.DeleteVerifySession(sessionID); err != nil {
		return nil, err
	}

	account, err := impl.db.UpdateState(id)
	if err != nil {
		return nil, err
	}

	return impl.newAccountSession(
		account.ID,
		account.AboType,
		verifyRequest.IPAddress,
		verifyRequest.UserAgent)
}

// tries to login a verified account
func (impl *Impl) Login(loginRequest *types.LoginAccountDTO) (*types.NewAccountSessionDTO, error) {
	account, err := impl.db.GetUserByEmail(loginRequest.Email)
	if err != nil {
		return nil, err
	}

	if _, err := checkpassword(loginRequest.Password, account.PasswordHash); err != nil {
		return nil, err
	}

	if !isValidEmail(loginRequest.Email) {
		return nil, errors.New("invalid crtedentials")
	}

	return impl.newAccountSession(
		account.ID,
		account.AboType,
		loginRequest.IPAddress,
		loginRequest.UserAgent)
}

// tries to register a password reset
func (impl *Impl) NewPassword(email string) (*types.VerifySessionDTO, error) {
	if !isValidEmail(email) {
		return nil, errors.New("invalid credentials")
	}

	account, err := impl.db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	number, cookie, err := impl.session.SetVerifyData(account.ID)
	if err != nil {
		return nil, err
	}

	if err := impl.sendVerifyData(number, account.Email, account.Name); err != nil {
		return nil, err
	}

	return &types.VerifySessionDTO{Name: account.Name, Cookie: cookie}, nil
}

// tries to execute a registered password reset
func (impl *Impl) ChangePassword(changePassRequest *types.ChangePasswordDTO) (*types.NewAccountSessionDTO, error) {
	if !passwordConfirmed(changePassRequest.Password, changePassRequest.Confirmation) {
		return nil, errors.New("invalid credentials")
	}

	passwordHash, err := getHashedPassword(changePassRequest.Password)
	if err != nil {
		return nil, err
	}

	id, sessionID, number, err := impl.session.GetVerifyData(changePassRequest.Cookie)
	if err != nil {
		return nil, err
	}

	if number != changePassRequest.VerifyCode {
		return nil, errors.New("invalid credentials")
	}

	account, err := impl.db.ChangePassword(id, passwordHash)
	if err != nil {
		return nil, err
	}

	if err := impl.session.DeleteVerifySession(sessionID); err != nil {
		return nil, err
	}

	return impl.newAccountSession(
		account.ID,
		account.AboType,
		changePassRequest.IPAddress,
		changePassRequest.UserAgent)
}

// tries to refresh a refresh session of a user who has accsess to one of them
func (impl *Impl) RefreshSession(refresh *types.RefreshSessionDTO) (*types.NewAccountSessionDTO, error) {
	newCookie, err := impl.session.VerifyRefreshData(refresh.Cookie, refresh.IPAddress, refresh.IPAddress)
	if err != nil {
		return nil, err
	}

	id, role, err := impl.jwt.Verify(refresh.JWT)
	if err != nil {
		return nil, err
	}

	newJWT, err := impl.jwt.Sign(id, role)
	if err != nil {
		return nil, err
	}

	return &types.NewAccountSessionDTO{
		Access:  newJWT,
		Refresh: newCookie,
	}, nil
}

// tries to remove a session of a users account who access to it
func (impl *Impl) RemoveSession(remove *types.RefreshSessionDTO) error {
	if _, _, err := impl.jwt.Verify(remove.JWT); err != nil {
		return err
	}

	return impl.session.RemoveRefreshData(remove.Cookie, remove.IPAddress, remove.UserAgent)
}
