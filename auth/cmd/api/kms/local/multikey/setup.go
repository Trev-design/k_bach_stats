package multikey

import (
	"auth_server/cmd/api/kms/local/key"
	"errors"
	"math"
	"sync"
	"time"
)

type keyHandler struct {
	timestamp time.Time
	enclave   *key.Enclave
}

type Manager struct {
	mutex            sync.RWMutex
	keys             map[string]*keyHandler
	createdTimeStamp time.Time
	updatedTimestamp time.Time
	intervalTicker   *time.Ticker
	numberOfKeys     int
	intervalDuration time.Duration
	expiryDuration   time.Duration
}

func NewManager(duration, expiry time.Duration, numverOfKeys int) *Manager {
	return &Manager{
		mutex:            sync.RWMutex{},
		keys:             make(map[string]*keyHandler),
		createdTimeStamp: time.Now().UTC(),
		updatedTimestamp: time.Now().UTC(),
		intervalTicker:   time.NewTicker(duration),
		numberOfKeys:     numverOfKeys,
		intervalDuration: duration,
		expiryDuration:   expiry,
	}
}

func (manager *Manager) GetKey(key string, timestamp time.Time) ([]byte, error) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	if time.Since(timestamp) > manager.expiryDuration {
		return nil, errors.New("expired key")
	}

	keyset, ok := manager.keys[key]
	if !ok {
		return nil, errors.New("invalid key")
	}

	index := int(math.Abs(float64(timestamp.Second() - manager.createdTimeStamp.Second())))

	return keyset.enclave.GetKey(index)
}

func (manager *Manager) StopKeyManager() error {
	for _, keyset := range manager.keys {
		if err := keyset.enclave.DestroyKey(); err != nil {
			return err
		}
	}

	return nil
}

func (manager *Manager) ComputeRotateInterval() {
	for range manager.intervalTicker.C {
		if len(manager.keys) < manager.numberOfKeys {
			manager.addKey()
		} else {
			manager.updateKeys()
		}
	}
}

func (manager *Manager) addKey() {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
}

func (manager *Manager) updateKeys() {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for _, keyset := range manager.keys {
		if time.Since(keyset.timestamp) > manager.expiryDuration {
			keyset.updateKey()
			return
		}
	}
}

func (handler *keyHandler) updateKey() {

}
