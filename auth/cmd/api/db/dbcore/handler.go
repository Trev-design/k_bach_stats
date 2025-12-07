package dbcore

import (
	"auth_server/cmd/api/domain/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Creates a new user. Per default he is not authorized.
func (db *Database) AddUser(newAccount *types.NewAccountDM) (string, error) {
	var guid uuid.UUID
	conn := db.conn.Get()
	conn.waitgroup.Add(1)
	defer conn.waitgroup.Done()
	if err := conn.conn.Transaction(func(tx *gorm.DB) error {
		id, err := addAccount(tx)
		if err != nil {
			return err
		}

		if err = addUser(
			tx,
			&userCreds{
				name:         newAccount.Name,
				email:        newAccount.Email,
				passwordHash: newAccount.PasswordHash,
				accountID:    id,
			},
		); err != nil {
			return err
		}

		if err = addRole(tx, id); err != nil {
			return err
		}

		guid = id

		return nil
	}); err != nil {
		return "", err
	}

	return guid.String(), nil
}

// Fetches user from Database
func (db *Database) GetUser(id uuid.UUID) (*types.AccountDM, error) {
	conn := db.conn.Get()
	conn.waitgroup.Add(1)
	defer conn.waitgroup.Done()
	return getUser(conn.conn, id)
}

// Fetches user from Database based on his email.
// This is common when the user tries to sign in.
func (db *Database) GetUserByEmail(email string) (*types.AccountDM, error) {
	conn := db.conn.Get()
	conn.waitgroup.Add(1)
	defer conn.waitgroup.Done()

	row := conn.conn.Table("accounts AS a").
		Joins("LEFT JOIN users AS u ON a.id = u.account_id").
		Joins("LEFT JOIN roles AS r ON a.id = r.account_id").
		Select("a.id, u.name, u.email, u.password_hash, r.is_verified, r.abo_type").
		Where("u.email = ?", email).
		Row()

	return userByEmail(row)
}

// When the user is verified. His authorized status changes and he is a valid use with an abo plan default community.
func (db *Database) UpdateState(id uuid.UUID) (*types.AccountDM, error) {
	conn := db.conn.Get()
	conn.waitgroup.Add(1)
	defer conn.waitgroup.Done()

	if err := conn.conn.Model(&Role{}).
		Where("account_id = ?", id).
		Updates(map[string]any{"is_verified": true, "abo_type": "COMMUNITY"}).
		Error; err != nil {
		return nil, err
	}

	return getUser(conn.conn, id)
}

// Fallback: If the user has forgotten their password or just want to change it, the user can change it.
func (db *Database) ChangePassword(id uuid.UUID, passwordHash string) (*types.AccountDM, error) {
	conn := db.conn.Get()
	conn.waitgroup.Add(1)
	defer conn.waitgroup.Done()

	if err := conn.conn.Model(&User{}).Where("account_id = ?", id).Update("password_hash", passwordHash).Error; err != nil {
		return nil, err
	}

	updatedAccount, err := getUser(conn.conn, id)
	if err != nil {
		return nil, err
	}

	return updatedAccount, err
}
