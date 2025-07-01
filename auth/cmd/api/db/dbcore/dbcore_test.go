package dbcore_test

import (
	"auth_server/cmd/api/db/dbcore"
	"auth_server/cmd/api/domain/types"
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type dbCreds struct {
	user     string
	password string
	database string
	host     string
	port     string
}

var creds *dbCreds
var db *dbcore.Database

func TestMain(m *testing.M) {
	newCreds, cancel := setupTestContainer()
	creds = new(dbCreds)
	creds = newCreds

	newDB, err := dbcore.NewDatabaseBuilder().
		User(creds.user).
		Password(creds.password).
		Host(creds.host).
		Port(creds.port).
		DBName(creds.database).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	db = newDB

	time.Sleep(1 * time.Second)

	code := m.Run()

	cancel()
	db.CloseDatabase()
	os.Exit(code)
}

func TestConnectionFailedFalseHost(t *testing.T) {
	_, err := dbcore.NewDatabaseBuilder().
		User(creds.user).
		Password(creds.password).
		Host("invalid host").
		Port(creds.port).
		DBName(creds.database).
		Build()
	if err == nil {
		t.Fatal("should get an error but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func TestConnectionFailedFalsePort(t *testing.T) {
	_, err := dbcore.NewDatabaseBuilder().
		User(creds.user).
		Password(creds.password).
		Host(creds.host).
		Port("123").
		DBName(creds.database).
		Build()
	if err == nil {
		t.Fatal("should get an error but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func TestConnectionFailedFalseUser(t *testing.T) {
	_, err := dbcore.NewDatabaseBuilder().
		User("invalid_user").
		Password(creds.password).
		Host(creds.host).
		Port(creds.port).
		DBName(creds.database).
		Build()
	if err == nil {
		t.Fatal("should get an error but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func TestConnectionFailedFalsePassword(t *testing.T) {
	_, err := dbcore.NewDatabaseBuilder().
		User(creds.user).
		Password("invalid_password").
		Host(creds.host).
		Port(creds.port).
		DBName(creds.database).
		Build()
	if err == nil {
		t.Fatal("should get an error but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func TestConnectionFailedFalseDBName(t *testing.T) {
	_, err := dbcore.NewDatabaseBuilder().
		User(creds.user).
		Password(creds.password).
		Host(creds.host).
		Port(creds.port).
		DBName("invalid_db_name").
		Build()
	if err == nil {
		t.Fatal("should get an error but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func TestConnectionSuccess(t *testing.T) {
	newDB, err := dbcore.NewDatabaseBuilder().
		User(creds.user).
		Password(creds.password).
		Host(creds.host).
		Port(creds.port).
		DBName(creds.database).
		Build()
	if err != nil {
		t.Fatal(err)
	}

	if err := newDB.CloseDatabase(); err != nil {
		t.Fatal(err)
	}
}

func TestAddUser(t *testing.T) {
	id, err := db.AddUser(&types.NewAccountDM{
		Name:         "Gisela Alesig",
		Email:        "gisela.alesig@gigi.gi",
		PasswordHash: "securly hashed password i swear",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("successful created account %s", id)
}

func TestGetUser(t *testing.T) {
	id, err := db.AddUser(&types.NewAccountDM{
		Name:         "Gisela Alesig",
		Email:        "gisela2.alesig@gigi.gi",
		PasswordHash: "securly hashed password i swear",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("id is %s", id)

	guid, err := uuid.Parse(id)
	if err != nil {
		t.Fatal(err)
	}

	account, err := db.GetUser(guid)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(account)
}

func TestGetUserByEmail(t *testing.T) {
	account, err := db.GetUserByEmail("gisela.alesig@gigi.gi")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(account)
}

func TestUpdateState(t *testing.T) {
	id, err := db.AddUser(&types.NewAccountDM{
		Name:         "Gisela Alesig",
		Email:        "gisela222.alesig@gigi.gi",
		PasswordHash: "securly hashed password i swear",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("id is %s", id)

	guid, err := uuid.Parse(id)
	if err != nil {
		t.Fatal(err)
	}

	account, err := db.UpdateState(guid)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(account)
}

func TestChangePassword(t *testing.T) {
	id, err := db.AddUser(&types.NewAccountDM{
		Name:         "Gisela Alesig",
		Email:        "gisela22.alesig@gigi.gi",
		PasswordHash: "securly hashed password i swear",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("id is %s", id)

	guid, err := uuid.Parse(id)
	if err != nil {
		t.Fatal(err)
	}

	account, err := db.ChangePassword(guid, "another new securely hashed password")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(account)
}

func setupTestContainer() (*dbCreds, func()) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(100 * time.Second),
	}

	container, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		log.Fatal(err)
	}

	host, err := container.Host(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	mappedPort, err := container.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(15 * time.Second)

	return &dbCreds{
		user:     "test_user",
		password: "test_password",
		database: "testdb",
		host:     host,
		port:     mappedPort.Port(),
	}, func() { container.Terminate(context.Background()) }
}
