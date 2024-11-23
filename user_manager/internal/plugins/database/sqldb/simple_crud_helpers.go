package sqldb

import (
	"database/sql"
	"user_manager/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func (db *DBAdapter) insertUser(payload *types.UserMessagePayload) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	defer handleTransaction(transaction, err)

	statement, err := transaction.Prepare("INSERT INTO users (id, entity) VALUES (UNHEX(REPLACE(?, '-', '')), ?);")
	if err != nil {
		return err
	}
	defer statement.Close()

	guid := uuid.New().String()

	_, err = statement.Exec(guid, payload.Entity)
	if err != nil {
		return err
	}

	return db.insertProfile(transaction, guid, payload)
}

func (db *DBAdapter) insertProfile(transaction *sql.Tx, userID string, payload *types.UserMessagePayload) error {
	statement, err := transaction.Prepare("INSERT INTO profiles (id, bio, user_id) VALUES (UNHEX(REPLACE(?, '-', '')), ?, UNHEX(REPLACE(?, '-', '')));")
	if err != nil {
		return err
	}
	defer statement.Close()

	guid := uuid.New().String()

	_, err = statement.Exec(guid, "", userID)
	if err != nil {
		return err
	}

	return db.insertContact(transaction, guid, payload)
}

func (db *DBAdapter) insertContact(transaction *sql.Tx, profileID string, payload *types.UserMessagePayload) error {
	statement, err := transaction.Prepare("INSERT INTO contacts (id, name, email, image_file_path, profile_id) VALUES (UNHEX(REPLACE(?, '-', '')), ?, ?, ?, UNHEX(REPLACE(?, '-', '')));")
	if err != nil {
		return err
	}
	defer statement.Close()

	guid := uuid.New().String()

	_, err = statement.Exec(guid, payload.Username, payload.Email, "", profileID)
	if err != nil {
		return err
	}

	return nil
}

func (db *DBAdapter) removeUser(entity string) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	defer handleTransaction(transaction, err)

	statement, err := transaction.Prepare("DELETE FROM users WHERE entity = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(entity); err != nil {
		return err
	}

	return nil
}
