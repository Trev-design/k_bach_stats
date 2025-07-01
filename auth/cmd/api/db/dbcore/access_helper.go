package dbcore

import (
	"auth_server/cmd/api/domain/types"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userCreds struct {
	name, email, passwordHash string
	accountID                 uuid.UUID
}

func addAccount(transaction *gorm.DB) (uuid.UUID, error) {
	account := new(Account)
	if err := transaction.Create(account).Error; err != nil {
		return uuid.Nil, err
	}

	return account.ID, nil
}

func addUser(transaction *gorm.DB, creds *userCreds) error {
	return transaction.Create(&User{
		Name:         creds.name,
		Email:        creds.email,
		PasswordHash: creds.passwordHash,
		AccountID:    creds.accountID,
	}).Error
}

func addRole(transaction *gorm.DB, accountID uuid.UUID) error {
	return transaction.Create(&Role{AccountID: accountID}).Error
}

func getUser(db *gorm.DB, id uuid.UUID) (*types.AccountDM, error) {
	account := new(Account)
	err := db.Preload("User").Preload("Role").First(account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &types.AccountDM{
		ID:           account.ID.String(),
		Name:         account.User.Name,
		Email:        account.User.Email,
		PasswordHash: account.User.PasswordHash,
		Isverified:   account.Role.IsVerified,
		AboType:      account.Role.AboType,
	}, nil
}

func userByEmail(row *sql.Row) (*types.AccountDM, error) {
	var guid uuid.NullUUID
	var name sql.NullString
	var email sql.NullString
	var passwordHash sql.NullString
	var isVerified sql.NullBool
	var aboType sql.NullString

	if err := row.Scan(
		&guid,
		&name,
		&email,
		&passwordHash,
		&isVerified,
		&aboType,
	); err != nil {
		return nil, err
	}

	if !guid.Valid || !name.Valid || !email.Valid ||
		!passwordHash.Valid || !isVerified.Valid || !aboType.Valid {
		return nil, errors.New("invalid payload")
	}

	return &types.AccountDM{
		ID:           guid.UUID.String(),
		Name:         name.String,
		Email:        email.String,
		PasswordHash: passwordHash.String,
		Isverified:   isVerified.Bool,
		AboType:      aboType.String,
	}, nil
}
