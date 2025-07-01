package jwtcore

import (
	"crypto/ed25519"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (service *JWTService) Sign(id, role string) (string, error) {
	timestamp := time.Now()
	claims := jwt.MapClaims{
		"sub":  id,
		"iss":  service.guid,
		"iat":  timestamp.UTC().Unix(),
		"exp":  timestamp.UTC().Add(15 * time.Minute).Unix(),
		"role": role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	token.Header["kid"] = timestamp.Format(time.RFC3339Nano)

	seed, err := service.seeds.GetSeed(timestamp)
	if err != nil {
		return "", err
	}

	key := ed25519.NewKeyFromSeed(seed)

	return token.SignedString(key)
}

func (service *JWTService) Verify(jwtToken string) (string, string, error) {
	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("invalid session credentials")
		}

		timestamp, err := time.Parse(time.RFC3339Nano, kid)
		if err != nil {
			return nil, err
		}

		seed, err := service.seeds.GetSeed(timestamp)
		if err != nil {
			return nil, err
		}

		privKey := ed25519.NewKeyFromSeed(seed)
		return privKey.Public(), nil
	})
	if err != nil {
		return "", "", err
	}

	if !token.Valid {
		return "", "", errors.New("invalid session credentials")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid session credentials")
	}

	id, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("invalid session credentials")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", errors.New("invalid session credentials")
	}

	return id, role, nil
}
