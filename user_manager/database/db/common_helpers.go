package db

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func (db *Database) insertItem(query string, args ...any) error {
	if _, err := db.Exec(query, args...); err != nil {
		return err
	}
	return nil
}

func (db *Database) removeItem(query string, args ...any) error {
	if _, err := db.Exec(query, args...); err != nil {
		return err
	}
	return nil
}

func validvalues(strings ...sql.NullString) bool {
	for _, val := range strings {
		if !val.Valid {
			return false
		}
	}

	return true
}

func (db *Database) insertUserCredentials(kind, foreignKey string, user *UserPayload) (string, error) {
	guid := uuid.New().String()

	switch kind {
	case "user":
		return guid, db.insertItem(insertNewUser, guid, user.Entity)
	case "profile":
		return guid, db.insertItem(insertNewProfile, guid, "", foreignKey)
	case "contact":
		return guid, db.insertItem(insertNewContact, guid, user.Username, user.Email, "", foreignKey)
	default:
		return "", errors.New("invalid insert option")
	}
}

func uuidFromBin(binaryIDs ...*string) error {
	for _, id := range binaryIDs {
		guid, err := uuid.FromBytes([]byte(*id))
		if err != nil {
			return err
		}

		*id = guid.String()
	}

	return nil
}

func handleTransaction(transaction *sql.Tx, err error) {
	if p := recover(); p != nil {
		transaction.Rollback()
	} else if err != nil {
		transaction.Rollback()
	} else if commitError := transaction.Commit(); commitError != nil {
		transaction.Rollback()
	}
}
