package impl

import (
	"auth_server/cmd/api/domain/types"
	"encoding/json"
	"errors"
	"regexp"
	"unicode"

	"github.com/alexedwards/argon2id"
)

func (impl *Impl) sendVerifyData(number, email, name string) error {
	mail := struct {
		Verify string `json:"verify"`
		Email  string `json:"email"`
		Name   string `json:"name"`
	}{
		Verify: number,
		Email:  email,
		Name:   name,
	}

	message, err := json.Marshal(&mail)
	if err != nil {
		return errors.New("something went wrong")
	}

	return impl.broker.SendMessage("verify_mail", message)
}

func (impl *Impl) newAccountSession(account, role, ip, userAgent string) (*types.NewAccountSessionDTO, error) {
	refresh, err := impl.session.SetRefreshData(account, ip, userAgent)
	if err != nil {
		return nil, err
	}

	access, err := impl.jwt.Sign(account, role)
	if err != nil {
		return nil, err
	}

	return &types.NewAccountSessionDTO{Access: access, Refresh: refresh}, nil
}

func isValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func passwordConfirmed(password, confirmation string) bool {
	return password == confirmation
}

func getHashedPassword(password string) (string, error) {
	if !validPassword(password) {
		return "", errors.New("invalid credentials")
	}

	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func checkpassword(password, passwordHash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, passwordHash)
}

func validPassword(password string) bool {
	hasCapitals := false
	hasLetters := false
	hasDigits := false
	hasSpecials := false

	if len(password) < 8 {
		return false
	}

	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasDigits = true

		case unicode.IsLower(char):
			hasLetters = true

		case unicode.IsUpper(char):
			hasCapitals = true

		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecials = true

		default:
			return false
		}
	}

	return hasCapitals && hasDigits && hasLetters && hasSpecials
}
