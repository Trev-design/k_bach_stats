package sessionkeys

func (keys *KeyManager) destroyKeys() error {
	if err := keys.oldKeys.DestroyKey(); err != nil {
		return err
	}

	return keys.newKeys.DestroyKey()
}
