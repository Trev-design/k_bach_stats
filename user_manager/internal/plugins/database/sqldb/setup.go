package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"user_manager/internal/plugins/database"

	_ "github.com/go-sql-driver/mysql"
)

func openDatabase() (*sql.DB, error) {
	return connect("IAmTheUser:ThisIsMyPassword@tcp(127.0.0.1:3306)/user_database")
}

func resetDatabase() (*sql.DB, error) {
	db, err := connect("IAmTheUser:ThisIsMyPassword@tcp(127.0.0.1:3306)/user_database")
	if err != nil {
		return nil, fmt.Errorf("line 21: %s", err.Error())
	}

	_, err = db.Exec("DROP DATABASE user_database;")
	if err != nil {
		return nil, database.ErrInternal
	}
	log.Println("database droped")

	db.Close()

	return initDatabase()
}

func initDatabase() (*sql.DB, error) {
	db, err := connect("IAmTheUser:ThisIsMyPassword@tcp(127.0.0.1:3306)/")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE DATABASE user_database;")
	if err != nil {
		return nil, err
	}
	log.Println("created database")

	db.Close()

	return initTables()
}

func initTables() (*sql.DB, error) {
	queries, err := getQueries()
	if err != nil {
		return nil, fmt.Errorf("line 55: %s", err.Error())
	}

	db, err := connect("IAmTheUser:ThisIsMyPassword@tcp(127.0.0.1:3306)/user_database")
	if err != nil {
		return nil, err
	}

	for _, query := range queries {
		if query == "" {
			continue
		}
		if _, err := db.Exec(query); err != nil {
			log.Println(err.Error())
			return nil, database.ErrForbiddenDatabaseRequest
		}
	}

	return db, nil
}

func connect(dsn string) (*sql.DB, error) {
	log.Println("connect")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, database.ErrInternal
	}
	log.Println("try ping")
	log.Println(err)

	if err = db.Ping(); err != nil {
		log.Println(err)
		return nil, database.ErrConnectionLost
	}

	return db, nil
}

func getQueries() ([]string, error) {
	queryFile, err := os.ReadFile("./init/mysql/init.sql")
	if err != nil {
		log.Println(err)
		return nil, database.ErrInternal
	}

	return strings.Split(string(queryFile), ";"), nil
}
