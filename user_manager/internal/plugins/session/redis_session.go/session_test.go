package redis_session

import (
	"encoding/json"
	"testing"
	"user_manager/graph/model"
	"user_manager/types"
)

// a test which tests a expected successful storage roundtrip
func Test_RedisUserManagementRoundTrip_success(test *testing.T) {
	// open redis store
	client, err := NewSessionAdapter()
	if err != nil {
		test.Error("could not start session store")
	}
	defer client.Close()

	// adding user
	session := &types.SessionMessagePayload{
		Name:      "MyveryCoolname",
		Account:   "123",
		SessionID: "321",
		AboType:   "TheSupaDupaPremiumAbo",
	}

	payload, err := json.Marshal(session)
	if err != nil {
		test.Error("could not make a json payload")
	}

	err = client.AddUser(payload)
	if err != nil {
		test.Errorf("could not put the payload in storage because %s", err.Error())
	}

	// remove user
	err = client.RemoveUser(payload)
	if err != nil {
		test.Errorf("could not remove payload from storage because %s", err.Error())
	}
}

// with an unsuccessfull add
func Test_RedisUserManagementRoundTrip_failed_add(test *testing.T) {
	// open redis store
	client, err := NewSessionAdapter()
	if err != nil {
		test.Error("could not start session store")
	}
	defer client.Close()

	// add user failure
	invalid := &model.UserEntity{
		User: "popel",
	}
	invalidPayload, err := json.Marshal(invalid)
	if err != nil {
		test.Error("could not make a json payload")
	}

	err = client.AddUser(invalidPayload)
	if err == nil {
		test.Errorf("expected an error")
	}

	// add user success
	session := &types.SessionMessagePayload{
		Name:      "MyveryCoolname",
		Account:   "123",
		SessionID: "321",
		AboType:   "TheSupaDupaPremiumAbo",
	}

	payload, err := json.Marshal(session)
	if err != nil {
		test.Error("could not make a json payload")
	}

	err = client.AddUser(payload)
	if err != nil {
		test.Errorf("could not put the payload in storage because %s", err.Error())
	}

	// remove user
	err = client.RemoveUser(payload)
	if err != nil {
		test.Errorf("could not remove payload from storage because %s", err.Error())
	}
}
