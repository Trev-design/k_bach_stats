package sessionkeys

import (
	"errors"
	"math"
	"time"
)

func (keys *KeyManager) GetKey(timeStamp time.Time) ([]byte, error) {
	keys.mutex.RLock()
	defer keys.mutex.RUnlock()

	if time.Since(timeStamp) > keys.duration {
		return nil, errors.New("expired key")
	}

	index := int(math.Abs(float64(timeStamp.Second() - keys.createdTimeStamp.Second())))

	if timeStamp.Before(keys.updatedTimeStamp) {
		return keys.oldKeys.GetKey(index)
	}

	return keys.newKeys.GetKey(index)
}
