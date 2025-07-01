package sessionenclave

func (key *Key) GetKey(diff int) ([]byte, error) {
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

func (key *Key) ChangeKey() error {
	key.mutex.Lock()
	defer key.mutex.Unlock()

	keyBuffer, err := key.enclave.Open()
	if err != nil {
		return err
	}

	keyBuffer.Melt()
	keyBuffer.Scramble()
	keyBuffer.Freeze()

	newEnclave := keyBuffer.Seal()
	key.enclave = newEnclave

	return nil
}
