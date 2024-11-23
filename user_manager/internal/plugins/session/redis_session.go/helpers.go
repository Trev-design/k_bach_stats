package redis_session

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"
	"user_manager/types"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidSession = errors.New("invalid session")

func (sessiona *SessionAdapter) parseClaims(token string) (*Claims, error) {
	claims := new(Claims)

	if _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("not able to authorize")
		}

		return sessiona.Secret, nil
	}); err != nil {
		return nil, err
	}

	return claims, nil
}

func (sessiona *SessionAdapter) getSession(key string) (*types.SessionMessagePayload, error) {
	payload, err := sessiona.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	sessionPayload := new(types.SessionMessagePayload)
	err = json.Unmarshal([]byte(payload), sessionPayload)
	if err != nil {
		return nil, err
	}

	return sessionPayload, nil
}

func checkSessionCredentials(session *types.SessionMessagePayload, claims *Claims) error {
	if !validSession(session, claims) {
		return ErrInvalidSession
	}

	return nil
}

func validSession(session *types.SessionMessagePayload, claims *Claims) bool {
	return session.Account == claims.Account && session.Name == claims.Name
}

func (sessiona *SessionAdapter) setNewUserSession(credentials *types.SessionMessagePayload, payload []byte) error {

	if status, err := sessiona.SetEx(
		context.Background(),
		credentials.SessionID,
		string(payload),
		(60*60*24)*time.Second,
	).Result(); err != nil {
		log.Println(status)
		return err
	}

	return nil
}
