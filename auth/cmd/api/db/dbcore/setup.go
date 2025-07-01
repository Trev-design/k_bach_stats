package dbcore

import (
	"auth_server/cmd/api/tlsconf"

	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

type DatabaseBuilder struct {
	user       string
	password   string
	host       string
	port       string
	dbname     string
	tlsBuilder *tlsconf.TLSBuilder
}

func NewDatabaseBuilder() *DatabaseBuilder {
	return &DatabaseBuilder{}
}

func (builder *DatabaseBuilder) User(user string) *DatabaseBuilder {
	builder.user = user
	return builder
}

func (builder *DatabaseBuilder) Password(password string) *DatabaseBuilder {
	builder.password = password
	return builder
}

func (builder *DatabaseBuilder) Host(host string) *DatabaseBuilder {
	builder.host = host
	return builder
}

func (builder *DatabaseBuilder) Port(port string) *DatabaseBuilder {
	builder.port = port
	return builder
}

func (builder *DatabaseBuilder) DBName(dbname string) *DatabaseBuilder {
	builder.dbname = dbname
	return builder
}

func (builder *DatabaseBuilder) WithTLS(tlsBuilder *tlsconf.TLSBuilder) *DatabaseBuilder {
	builder.tlsBuilder = tlsBuilder
	return builder
}

func (builder *DatabaseBuilder) Build() (*Database, error) {
	db, err := builder.makeConn()
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Account{}, &User{}, &Role{}); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (db *Database) CloseDatabase() error {
	sqlDB, err := db.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
