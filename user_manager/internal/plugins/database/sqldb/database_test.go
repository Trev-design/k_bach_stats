package sqldb_test

import (
	"encoding/json"
	"testing"
	"user_manager/internal/plugins/database/sqldb"
	"user_manager/types"
)

func Test_SQLUserManagerInsert_success(test *testing.T) {
	client, err := sqldb.NewDatabaseAdapter("keep")
	if err != nil {
		test.Error(err.Error())
	}
	defer client.Close()

	user := &types.UserMessagePayload{
		Entity:   "135",
		Username: "johnny",
		Email:    "some@email.com",
	}

	payload, err := json.Marshal(user)
	if err != nil {
		test.Error("something went wrong")
	}

	if err := client.AddUser(payload); err != nil {
		test.Error("could not push to database")
	}
}

func Test_DatabaseSelectUserByID_success(test *testing.T) {
	client, err := sqldb.NewDatabaseAdapter("keep")
	if err != nil {
		test.Error("could not start database")
	}
	defer client.Close()

	userEntity, err := client.UserID("135")
	if err != nil {
		test.Error("no user with this entity")
	}

	if userEntity.User == "" || userEntity == nil {
		test.Error("invalid entity")
	}

	_, err = client.UserByID(userEntity.User)
	if err != nil {
		test.Error("user not found")
	}
}

func Test_SQLUserManagerRemove_success(test *testing.T) {
	client, err := sqldb.NewDatabaseAdapter("keep")
	if err != nil {
		test.Error("could not start database")
	}
	defer client.Close()

	payload := "135"
	if err := client.RemoveUser([]byte(payload)); err != nil {
		test.Error("could not remove user from database")
	}
}
