package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func resetDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", "")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(dropDatabase); err != nil {
		return nil, err
	}
	db.Close()

	return initDatabase()
}

func initDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", "IAmTheUser:ThisIsMyPassword@tcp(127.0.0.1:3306)/")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createDatatabas); err != nil {
		return nil, err
	}

	db.Close()

	return initTables()
}

func initTables() (*sql.DB, error) {
	db, err := sql.Open("mysql", "")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	initQueryFile, err := os.ReadFile("./sql/init.sql")
	if err != nil {
		log.Printf("could not find file %v\n", err)
		return nil, err
	}

	queries := strings.Split(string(initQueryFile), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		fmt.Println(query)
		if query == "" {
			continue
		}

		fmt.Println(query)

		if _, err := db.Exec(query); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}

	return db, nil
}
