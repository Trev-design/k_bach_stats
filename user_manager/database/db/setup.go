package db

import (
	"database/sql"
	"errors"
)

func Setup(opt string) (*Database, error) {
	db, err := computeSetup(opt)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func computeSetup(opt string) (*sql.DB, error) {
	switch opt {
	case "reset":
		return resetDatabase()
	case "create":
		return initDatabase()
	case "keep":
		return sql.Open("mysql", "IAmTheUser:ThisIsMyPassword@tcp(127.0.0.1:3306)/user_database")
	default:
		return nil, errors.New("invalid option")
	}
}
