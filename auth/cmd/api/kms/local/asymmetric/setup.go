package asymmetric

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"sync"

	"github.com/awnumar/memguard"
)

type KeyManager struct {
	privateKey *memguard.Enclave
	publicKey  string
	mutex      sync.Mutex
}

func NewKeyManager() (*KeyManager, error) {
	priv, pub, err := makeKeys()
	if err != nil {
		return nil, err
	}

	return &KeyManager{
		privateKey: memguard.NewEnclave(priv),
		publicKey:  base64.RawURLEncoding.EncodeToString(pub),
		mutex:      sync.Mutex{},
	}, nil
}

func makeKeys() ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	pub := &priv.PublicKey

	privateBytes := x509.MarshalPKCS1PrivateKey(priv)
	publicBytes := x509.MarshalPKCS1PublicKey(pub)

	return privateBytes, publicBytes, nil
}

func (keys *KeyManager) Decrypt(cipher []byte) ([]byte, error) {
	keys.mutex.Lock()
	defer keys.mutex.Unlock()

	keybuffer, err := keys.privateKey.Open()
	if err != nil {
		return nil, err
	}

	keyBytes := keybuffer.Bytes()
	key, err := x509.ParsePKCS1PublicKey(keyBytes)
	if err != nil {
		return nil, err
	}

	keys.privateKey = keybuffer.Seal()

	return rsa.EncryptOAEP(sha256.New(), rand.Reader, key, cipher, nil)
}

func (keys *KeyManager) SwapAndGet() (string, error) {
	keys.mutex.Lock()
	defer keys.mutex.Unlock()

	priv, pub, err := makeKeys()
	if err != nil {
		return "", err
	}

	keyBuffer, err := keys.privateKey.Open()
	if err != nil {
		return "", err
	}

	keyBuffer.Melt()

	for index, char := range priv {
		keyBuffer.Bytes()[index] = char
	}

	keyBuffer.Freeze()

	publicKey := base64.RawURLEncoding.EncodeToString(pub)

	keys.privateKey = keyBuffer.Seal()
	keys.publicKey = publicKey

	return publicKey, err
}
