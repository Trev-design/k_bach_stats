package key

import (
	"errors"
	"sync"

	"github.com/awnumar/memguard"
)

type Enclave struct {
	mutex        sync.Mutex
	container    *memguard.Enclave
	numberOfKeys int
}

func NewEnclave(numberOfKeys int, keyBytes []byte) *Enclave {
	return &Enclave{
		mutex:        sync.Mutex{},
		container:    memguard.NewEnclave(keyBytes),
		numberOfKeys: numberOfKeys,
	}
}

func (enclave *Enclave) GetKey(diff int) ([]byte, error) {
	if enclave.container == nil {
		return nil, errors.New("enclave already destroyed")
	}

	enclave.mutex.Lock()
	defer enclave.mutex.Unlock()

	index := diff % enclave.numberOfKeys

	keyBuffer, err := enclave.container.Open()
	if err != nil {
		return nil, err
	}
	keyChunk := keyBuffer.Bytes()

	keyBytes := keyChunk[index*32 : (index+1)*32]

	keyData := make([]byte, len(keyBytes))

	copy(keyData, keyBytes)

	newContainer := keyBuffer.Seal()
	enclave.container = newContainer

	return keyData, nil
}

func (enclave *Enclave) DestroyKey() error {
	enclave.mutex.Lock()
	defer enclave.mutex.Unlock()

	keyBuffer, err := enclave.container.Open()
	if err != nil {
		return err
	}

	keyBuffer.Destroy()
	enclave.container = nil

	return nil
}

func (enclave *Enclave) ChangeKey(newKeyBytes []byte) error {
	enclave.mutex.Lock()
	defer enclave.mutex.Unlock()

	keyBuffer, err := enclave.container.Open()
	if err != nil {
		return err
	}

	keyBuffer.Melt()

	keyBytes := keyBuffer.Bytes()

	for index := range keyBytes {
		keyBytes[index] = newKeyBytes[index]
	}

	keyBuffer.Freeze()

	newContainer := keyBuffer.Seal()
	enclave.container = newContainer

	return nil
}
