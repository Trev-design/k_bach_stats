package singlekey

import (
	"auth_server/cmd/api/kms/local/key"
	"sync"
	"time"
)

type Manager struct {
	mutex            sync.Mutex
	newKeys          *key.Enclave
	oldKeys          *key.Enclave
	updatedTimestamp time.Time
	createdTimestamp time.Time
	intervalTicker   *time.Ticker
	duration         time.Duration
}

func NewManager(duration time.Duration) *Manager {
	return &Manager{
		mutex:            sync.Mutex{},
		newKeys:          key.NewEnclave(12, []byte("")),
		oldKeys:          key.NewEnclave(12, []byte("")),
		updatedTimestamp: time.Now().UTC(),
		createdTimestamp: time.Now().UTC(),
		intervalTicker:   time.NewTicker(duration),
		duration:         duration,
	}
}
