package sqldb

import (
	"encoding/json"
	"user_manager/types"
)

func (dba *DBAdapter) AddUser(payload []byte) error {
	newUser := new(types.UserMessagePayload)
	if err := json.Unmarshal(
		payload,
		newUser,
	); err != nil {
		return err
	}

	return dba.insertUser(newUser)
}

func (dba *DBAdapter) RemoveUser(payload []byte) error {
	return dba.removeUser(string(payload))
}
