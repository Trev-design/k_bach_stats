package singlekey

import (
	"auth_server/cmd/api/kms/local/key"
	"errors"
	"log"
	"math"
	"sync"
	"time"
)

type Manager struct {
	mutex            sync.RWMutex
	newKeys          *key.Enclave
	oldKeys          *key.Enclave
	updatedTimestamp time.Time
	createdTimestamp time.Time
	intervalTicker   *time.Ticker
	duration         time.Duration
}

func NewManager(duration time.Duration) *Manager {
	return &Manager{
		mutex:            sync.RWMutex{},
		newKeys:          key.NewEnclave(12, []byte("")),
		oldKeys:          key.NewEnclave(12, make([]byte, 12*32)),
		updatedTimestamp: time.Now().UTC(),
		createdTimestamp: time.Now().UTC(),
		intervalTicker:   time.NewTicker(duration),
		duration:         duration,
	}
}

func (manager *Manager) StopKeyManager() error {
	manager.intervalTicker.Stop()

	return manager.destroyKeys()
}

func (manager *Manager) ComputeRotateInterval() {
	for range manager.intervalTicker.C {
		manager.mutex.Lock()

		manager.oldKeys, manager.newKeys = manager.newKeys, manager.oldKeys

		if err := manager.newKeys.ChangeKey([]byte("")); err != nil {
			log.Println(err)
		}

		manager.updatedTimestamp = time.Now().UTC()

		manager.mutex.Unlock()
	}
}

func (manager *Manager) GetKey(timestamp time.Time) ([]byte, error) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	if time.Since(timestamp) > manager.duration {
		return nil, errors.New("expired key")
	}

	index := int(math.Abs(float64(timestamp.Second() - manager.createdTimestamp.Second())))

	if timestamp.Before(manager.updatedTimestamp) {
		return manager.oldKeys.GetKey(index)
	}

	return manager.newKeys.GetKey(index)
}

func (manager *Manager) destroyKeys() error {
	if err := manager.oldKeys.DestroyKey(); err != nil {
		return err
	}

	return manager.newKeys.DestroyKey()
}
