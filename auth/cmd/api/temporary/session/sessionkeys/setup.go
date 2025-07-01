package sessionkeys

import (
	"auth_server/cmd/api/temporary/session/sessionenclave"
	"log"
	"sync"
	"time"
)

type KeyManager struct {
	mutex            sync.RWMutex
	newKeys          *sessionenclave.Key
	oldKeys          *sessionenclave.Key
	updatedTimeStamp time.Time
	createdTimeStamp time.Time
	intervalTicker   *time.Ticker
	duration         time.Duration
}

func NewKeyManager(duration time.Duration) *KeyManager {
	return &KeyManager{
		mutex:            sync.RWMutex{},
		newKeys:          sessionenclave.NewKey(12),
		oldKeys:          sessionenclave.NewKey(12),
		updatedTimeStamp: time.Now().UTC(),
		createdTimeStamp: time.Now().UTC(),
		intervalTicker:   time.NewTicker(duration),
		duration:         duration,
	}
}

func (keys *KeyManager) StopKeyManager() error {
	keys.intervalTicker.Stop()

	return keys.destroyKeys()
}

func (keys *KeyManager) ComputeRotateInterval() {
	for range keys.intervalTicker.C {
		keys.mutex.Lock()

		keys.oldKeys, keys.newKeys = keys.newKeys, keys.oldKeys

		if err := keys.newKeys.ChangeKey(); err != nil {
			log.Println(err)
		}

		keys.updatedTimeStamp = time.Now().UTC()

		keys.mutex.Unlock()
	}
}
