package redissession

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"
	"user_manager/middleware"

	"github.com/golang-jwt/jwt/v5"
	redis "github.com/redis/go-redis/v9"
)

const dayInSeconds = 60 * 60 * 24

type SessionClient struct {
	Secret []byte
	*redis.Client
}

type Session struct {
	Name    string `json:"name"`
	Account string `json:"account"`
	ID      string `json:"id"`
	AboType string `json:"abo_type"`
}

func Setup() (*SessionClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       2,
	})

	if ping, err := client.Ping(context.Background()).Result(); err != nil {
		log.Println(ping)
		return nil, err
	}

	key, err := os.ReadFile("public.pem")
	if err != nil {
		return nil, err
	}

	return &SessionClient{Client: client, Secret: key}, nil
}

func (client *SessionClient) AddUser(payload []byte) error {
	session := new(Session)
	if err := json.Unmarshal(payload, session); err != nil {
		return err
	}

	if status, err := client.SetEx(
		context.Background(),
		session.ID,
		string(payload),
		dayInSeconds*time.Second,
	).Result(); err != nil {
		log.Println(status)
		return err
	}

	return nil
}

func (client *SessionClient) RemoveUser(payload []byte) error {
	session := new(Session)
	if err := json.Unmarshal(payload, session); err != nil {
		return err
	}

	if err := client.Del(context.Background(), session.ID).Err(); err != nil {
		// TODO: Add Error Handler
		return err
	}

	return nil
}

func (client *SessionClient) CheckSession(token string) error {
	claims := new(middleware.Claims)

	if _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("false signing method")
		}
		return client.Secret, nil
	}); err != nil {
		log.Println(err.Error())
		return err
	}

	payload, err := client.Get(context.Background(), claims.Session).Result()
	if err != nil {
		return err
	}

	session := new(Session)

	if err := json.Unmarshal([]byte(payload), session); err != nil {
		return err
	}

	if session.Account != claims.Account {
		return errors.New("invalid account")
	}

	if session.ID != claims.ID {
		return errors.New("invalid ACCOUNT")
	}

	if session.Name != claims.Name {
		return errors.New("INVALID ACCOUNT")
	}

	return nil
}
