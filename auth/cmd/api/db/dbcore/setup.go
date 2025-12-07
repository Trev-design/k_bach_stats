package dbcore

import (
	"auth_server/cmd/api/tlsconf"
	"auth_server/cmd/api/utils/connection"
	"sync"

	"gorm.io/gorm"
)

// we use an orm called gorm to compute our database resources
type Database struct {
	credentialsChannel chan connection.Credentials
	builder            *DatabaseBuilder
	conn               *connection.Handler[Connection]
}

type Connection struct {
	waitgroup *sync.WaitGroup
	conn      *gorm.DB
}

type DatabaseBuilder struct {
	user       string
	password   string
	host       string
	port       string
	dbname     string
	tlsBuilder *tlsconf.TLSBuilder
	pipe       chan connection.Credentials
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

func (builder *DatabaseBuilder) WithCredentialsPipe(pipe chan connection.Credentials) *DatabaseBuilder {
	builder.pipe = pipe
	return builder
}

func (builder *DatabaseBuilder) Build() (*Database, error) {
	conn, err := connection.NewBuilder(builder).Build()
	if err != nil {
		return nil, err
	}

	return &Database{
		conn:               conn,
		credentialsChannel: builder.pipe,
		builder:            builder,
	}, nil
}

func (db *Database) Rotate() {
	for cred := range db.credentialsChannel {
		db.builder.User(cred.UserName).Password(cred.Password)
		db.conn.Rotate(db.builder)
	}
}

func (db *Database) CloseDatabase() error {
	return db.conn.Get().Close()
}
