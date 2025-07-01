package jwtenclave

import (
	"sync"

	"github.com/awnumar/memguard"
)

type Enclave struct {
	mutex   sync.Mutex
	enclave *memguard.Enclave
}

func NewEnclave() *Enclave {
	return &Enclave{
		enclave: memguard.NewEnclaveRandom(32),
		mutex:   sync.Mutex{},
	}
}
