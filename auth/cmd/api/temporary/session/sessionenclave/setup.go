package sessionenclave

import (
	"sync"

	"github.com/awnumar/memguard"
)

type Key struct {
	mutex   sync.Mutex
	enclave *memguard.Enclave
	numKeys int
}

func NewKey(numKeys int) *Key {
	return &Key{
		enclave: memguard.NewEnclaveRandom(32 * numKeys),
		mutex:   sync.Mutex{},
		numKeys: numKeys,
	}
}
