package jwtmanager

import (
	"errors"
	"time"
)

func (manager *SeedManager) GetSeed(timestamp time.Time) ([]byte, error) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	if time.Since(timestamp) > manager.intervalDuration {
		return nil, errors.New("expired key")
	}

	if timestamp.Before(manager.updatedTimestamp) {
		return manager.oldSeeds.GetSeed()
	}

	return manager.newSeeds.GetSeed()
}
