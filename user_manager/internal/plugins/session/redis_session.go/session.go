package redis_session

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"user_manager/types"

	"github.com/golang-jwt/jwt/v5"
	redis "github.com/redis/go-redis/v9"
)

type SessionAdapter struct {
	*redis.Client
	Secret *rsa.PublicKey
}

type Claims struct {
	Name    string `json:"name"`
	Account string `json:"entity"`
	AboType string `json:"abo"`
	Session string `json:"session_id"`
	jwt.RegisteredClaims
}

func NewSessionAdapter() (*SessionAdapter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       2,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	pemFile, err := os.ReadFile("./certs/public.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemFile)
	if block == nil {
		return nil, errors.New("could not decode pem key")
	}

	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &SessionAdapter{
		Client: client,
		Secret: key,
	}, nil
}

func (sessiona *SessionAdapter) AddUser(payload []byte) error {
	sessionCredentials := new(types.SessionMessagePayload)
	err := json.Unmarshal(payload, sessionCredentials)
	if err != nil {
		return err
	}

	return sessiona.setNewUserSession(sessionCredentials, payload)
}

func (sessiona *SessionAdapter) RemoveUser(payload []byte) error {
	sessionCredentials := new(types.SessionMessagePayload)
	err := json.Unmarshal(payload, sessionCredentials)
	if err != nil {
		return err
	}

	return sessiona.Del(context.Background(), sessionCredentials.SessionID).Err()
}

func (sessiona *SessionAdapter) CheckSesssion(token string) error {
	claims, err := sessiona.parseClaims(token)
	if err != nil {
		return err
	}

	sessionpayload, err := sessiona.getSession(claims.Session)
	if err != nil {
		return err
	}

	return checkSessionCredentials(sessionpayload, claims)
}

func (sessiona *SessionAdapter) InitialAuth(token string) (string, error) {
	claims, err := sessiona.parseClaims(token)
	if err != nil {
		return "", err
	}

	return claims.Account, nil
}
