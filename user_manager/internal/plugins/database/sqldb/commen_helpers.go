package sqldb

import (
	"database/sql"
	"user_manager/internal/plugins/database"

	"github.com/google/uuid"
)

func computeSetup(setupOpt string) (*sql.DB, error) {
	switch setupOpt {
	case "reset":
		return resetDatabase()

	case "keep":
		return openDatabase()

	default:
		return nil, database.ErrForbiddenInitOpt
	}
}

func uuidsFrombin(binaryIDs ...*string) error {
	for _, id := range binaryIDs {
		guid, err := uuidFromBin(*id)
		if err != nil {
			return err
		}

		*id = guid
	}

	return nil
}

func uuidFromBin(binaryID string) (string, error) {
	guid, err := uuid.FromBytes([]byte(binaryID))
	if err != nil {
		return "", database.ErrInvalidBinaryID
	}

	return guid.String(), nil
}

func validValues(valids ...bool) bool {
	for _, valid := range valids {
		if !valid {
			return false
		}
	}

	return true
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
