package main_test

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"testing"
	"time"
	"user_manager/graph/model"
	"user_manager/internal/plugins/session/redis_session.go"
	"user_manager/types"

	"github.com/golang-jwt/jwt/v5"
)

// a test which tests a expected successful storage roundtrip
func Test_RedisUserManagementRoundTrip_success(test *testing.T) {
	// open redis store
	client, err := redis_session.NewSessionAdapter()
	if err != nil {
		test.Error("could not start session store")
	}
	defer client.Close()

	// adding user
	session := &types.SessionMessagePayload{
		Name:      "MyveryCoolname",
		Account:   "123",
		SessionID: "321",
		AboType:   "TheSupaDupaPremiumAbo",
	}

	payload, err := json.Marshal(session)
	if err != nil {
		test.Error("could not make a json payload")
	}

	err = client.AddUser(payload)
	if err != nil {
		test.Errorf("could not put the payload in storage because %s", err.Error())
	}

	// remove user
	err = client.RemoveUser(payload)
	if err != nil {
		test.Errorf("could not remove payload from storage because %s", err.Error())
	}
}

// with an unsuccessfull add
func Test_RedisUserManagementRoundTrip_failed_add(test *testing.T) {
	// open redis store
	client, err := redis_session.NewSessionAdapter()
	if err != nil {
		test.Error("could not start session store")
	}
	defer client.Close()

	// add user failure
	invalid := &model.UserEntity{
		User: "popel",
	}
	invalidPayload, err := json.Marshal(invalid)
	if err != nil {
		test.Error("could not make a json payload")
	}

	err = client.AddUser(invalidPayload)
	if err == nil {
		test.Errorf("expected an error")
	}

	// add user success
	session := &types.SessionMessagePayload{
		Name:      "MyveryCoolname",
		Account:   "123",
		SessionID: "321",
		AboType:   "TheSupaDupaPremiumAbo",
	}

	payload, err := json.Marshal(session)
	if err != nil {
		test.Error("could not make a json payload")
	}

	err = client.AddUser(payload)
	if err != nil {
		test.Errorf("could not put the payload in storage because %s", err.Error())
	}

	// remove user
	err = client.RemoveUser(payload)
	if err != nil {
		test.Errorf("could not remove payload from storage because %s", err.Error())
	}
}

func Test_GenerateToken(test *testing.T) {
	key, err := makeKey()
	if err != nil {
		test.Error(err)
	}

	token, err := makeToken(key)
	if err != nil {
		test.Error(err)
	}

	log.Printf("%s\n", token)
}

func Test_CheckSession(test *testing.T) {
	token, err := makeCredentials()
	if err != nil {
		test.Error(err)
	}

	client, err := redis_session.NewSessionAdapter()
	if err != nil {
		test.Error(err)
	}
	defer client.Close()

	session := &types.SessionMessagePayload{
		Name:      "johnny",
		Account:   "135",
		SessionID: "321",
		AboType:   "TheSupaDupaPremiumAbo",
	}

	payload, err := json.Marshal(session)
	if err != nil {
		test.Error("could not make a json payload")
	}

	err = client.AddUser(payload)
	if err != nil {
		test.Errorf("could not put the payload in storage because %s", err.Error())
	}

	if err = client.CheckSesssion(token); err != nil {
		test.Error(err)
	}

	err = client.RemoveUser(payload)
	if err != nil {
		test.Errorf("could not remove payload from storage because %s", err.Error())
	}
}

func makeKey() (*rsa.PrivateKey, error) {
	pemData, err := os.ReadFile("./certs/private.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("no valid pemblock")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func makeToken(key *rsa.PrivateKey) (string, error) {
	ttl := 24 * time.Hour
	date := time.Now().UTC()

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{
			"exp":        date.Add(ttl).Unix(),
			"iat":        date.Unix(),
			"name":       "johnny",
			"entity":     "135",
			"abo":        "MySupaDupaPremiumAbo",
			"session_id": "321",
		},
	)

	return token.SignedString(key)
}

func makeCredentials() (string, error) {
	key, err := makeKey()
	if err != nil {
		return "", err
	}

	return makeToken(key)
}
