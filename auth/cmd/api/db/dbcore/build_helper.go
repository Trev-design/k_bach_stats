package dbcore

import (
	"errors"
	"fmt"
	"math"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (builder *DatabaseBuilder) makeConn() (*gorm.DB, error) {
	dsn, err := builder.makeDSN()
	if err != nil {
		return nil, err
	}

	var backoff time.Duration

	for counter := range 6 {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: true})
		if err == nil {
			if err := makeConnectionPool(db); err != nil {
				return nil, err
			}

			return db, nil
		} else {
			backoff = time.Duration(math.Pow(float64(counter), 2)) * time.Second
			time.Sleep(backoff)
		}
	}

	return nil, errors.New("something went wrong")
}

func (builder *DatabaseBuilder) makeDSN() (string, error) {
	if builder.tlsBuilder != nil {
		configMap, err := builder.tlsBuilder.BuildConfigMap()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=verify-full sslrootcert=%s sslcert=%s sslkey=%s",
			builder.user, builder.password, builder.host, builder.port, builder.dbname, configMap.CACertPath(), configMap.CertPath(), configMap.KeyPath(),
		), nil
	}

	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		builder.user, builder.password, builder.host, builder.port, builder.dbname,
	), nil
}

func makeConnectionPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return nil
}
