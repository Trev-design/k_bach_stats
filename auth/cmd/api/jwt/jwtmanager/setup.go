package jwtmanager

import (
	"auth_server/cmd/api/jwt/jwtenclave"
	"sync"
	"time"
)

type SeedManager struct {
	newSeeds         *jwtenclave.Enclave
	oldSeeds         *jwtenclave.Enclave
	createdTimestamp time.Time
	updatedTimestamp time.Time
	intervalTicker   *time.Ticker
	intervalDuration time.Duration
	mutex            sync.RWMutex
}

func NewSeedManager(duration time.Duration) *SeedManager {
	return &SeedManager{
		newSeeds:         jwtenclave.NewEnclave(),
		oldSeeds:         jwtenclave.NewEnclave(),
		createdTimestamp: time.Now().UTC(),
		updatedTimestamp: time.Now().UTC(),
		intervalTicker:   time.NewTicker(duration),
		intervalDuration: duration,
		mutex:            sync.RWMutex{},
	}
}

func (manager *SeedManager) CloseSeedManager() error {
	manager.intervalTicker.Stop()

	if err := manager.newSeeds.DestroySeeds(); err != nil {
		return err
	}

	return manager.oldSeeds.DestroySeeds()
}

func (manager *SeedManager) ComputeInterval() {
	for range manager.intervalTicker.C {
		manager.mutex.Lock()

		manager.newSeeds, manager.oldSeeds = manager.oldSeeds, manager.newSeeds
		manager.newSeeds.NewSeeds()

		manager.updatedTimestamp = time.Now().UTC()

		manager.mutex.Unlock()
	}
}
