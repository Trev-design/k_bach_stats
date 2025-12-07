package vaultclient

import (
	"context"
	"errors"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
)

func GetSecret() (string, error) {
	config := vault.DefaultConfig()

	client, err := vault.NewClient(config)
	if err != nil {
		return "", err
	}

	auth, err := auth.NewKubernetesAuth("auth")
	if err != nil {
		return "", err
	}

	info, err := client.Auth().Login(context.Background(), auth)
	if err != nil {
		return "", err
	}

	if info == nil {
		return "", errors.New("no auth info was returned after login")
	}

	secret, err := client.KVv2("secret").Get(context.Background(), "master")
	if err != nil {
		return "", err
	}

	value, ok := secret.Data["key"].(string)
	if !ok {
		return "", errors.New("value type assertion failed")
	}

	return value, nil
}
