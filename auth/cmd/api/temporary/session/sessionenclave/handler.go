package sessionenclave

import "errors"

// tries to fetch a key which is stored in memlock.
// gets a slice of bytes, which stores the fetched key bytes on success.
// on failure you'll get an error.
func (key *Key) GetKey(diff int) ([]byte, error) {
	if key.enclave == nil {
		return nil, errors.New("does not exist or already destoyed")
	}

	key.mutex.Lock()
	defer key.mutex.Unlock()

	index := diff % key.numKeys

	keyBuffer, err := key.enclave.Open()
	if err != nil {
		return nil, err
	}
	keyChunk := keyBuffer.Bytes()

	keyBytes := keyChunk[index*32 : (index+1)*32]

	keyData := make([]byte, len(keyBytes))

	copy(keyData, keyBytes)

	newEnclave := keyBuffer.Seal()
	key.enclave = newEnclave

	return keyData, nil
}

// deletes a set of keys out of the memlock.
func (key *Key) DestroyKey() error {
	key.mutex.Lock()
	defer key.mutex.Unlock()

	keyBuffer, err := key.enclave.Open()
	if err != nil {
		return err
	}

	keyBuffer.Destroy()
	key.enclave = nil

	return nil
}

// generates new keyset.
func (key *Key) ChangeKey() error {
	key.mutex.Lock()
	defer key.mutex.Unlock()

	keyBuffer, err := key.enclave.Open()
	if err != nil {
		return err
	}

	keyBuffer.Melt()     // must made mutable before make new keyset
	keyBuffer.Scramble() // generate new random bytes
	keyBuffer.Freeze()   // make key set imutable again

	newEnclave := keyBuffer.Seal()
	key.enclave = newEnclave

	return nil
}
