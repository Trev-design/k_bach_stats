package jwtenclave

import "errors"

// you get the an EdDSA seed from this function on success othervice an error
func (enclave *Enclave) GetSeed() ([]byte, error) {
	if enclave.enclave != nil {
		enclave.mutex.Lock()
		defer enclave.mutex.Unlock()

		buffer, err := enclave.enclave.Open()
		if err != nil {
			return nil, err
		}

		seedBytes := buffer.Bytes()
		seed := make([]byte, len(seedBytes))
		copy(seed, seedBytes)

		newEnclave := buffer.Seal()
		enclave.enclave = newEnclave

		return seed, nil
	}

	return nil, errors.New("enclave destroyed")
}

// changes the seeds on success othervice you'll get an error
func (enclave *Enclave) NewSeeds() error {
	if enclave.enclave != nil {
		enclave.mutex.Lock()
		defer enclave.mutex.Unlock()

		buffer, err := enclave.enclave.Open()
		if err != nil {
			return err
		}

		buffer.Melt()
		buffer.Scramble()
		buffer.Freeze()

		newEnclave := buffer.Seal()
		enclave.enclave = newEnclave

		return nil
	}

	return errors.New("enclave destroyed")
}

// destroys the seeds on success otherwice you'll get an error
func (enclave *Enclave) DestroySeeds() error {
	if enclave.enclave != nil {
		enclave.mutex.Lock()
		defer enclave.mutex.Unlock()

		buffer, err := enclave.enclave.Open()
		if err != nil {
			return err
		}

		buffer.Destroy()

		enclave.enclave = nil

		return nil
	}

	return errors.New("enclave destroyed")
}
