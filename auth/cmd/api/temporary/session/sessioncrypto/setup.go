package sessioncrypto

import (
	"auth_server/cmd/api/temporary/session/sessionkeys"
	"time"
)

type Crypt struct {
	keys *sessionkeys.KeyManager
}

func NewCrypt(duration time.Duration) (*Crypt, error) {
	return &Crypt{
		keys: sessionkeys.NewKeyManager(duration),
	}, nil
}

func (crypto *Crypt) CloseCrypto() error {
	return crypto.keys.StopKeyManager()
}

func (crypto *Crypt) ComputeRotateInterval() {
	crypto.keys.ComputeRotateInterval()
}
