package main_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"testing"
	"user_manager/graph/model"
	"user_manager/internal/plugins/database/sqldb"
	"user_manager/types"
)

func Test_InvitationRoundtrip(test *testing.T) {
	adapter, err := sqldb.NewDatabaseAdapter("keep")
	if err != nil {
		test.Error(err)
	}
	defer adapter.Close()

	// making user 1
	userID, contactID, entity, err := getEssentialInvitationData(
		adapter,
		"supa entity",
		"dupa username",
		"supa@dupa.mupa",
	)
	if err != nil {
		test.Error(err)
	}

	userID2, _, entity2, err := getEssentialInvitationData(
		adapter,
		"supa dupa entity",
		"dupa supa entity",
		"supadupa@dupasupa.mupa",
	)
	if err != nil {
		test.Error(err)
	}

	// make workspace
	workspaceID, err := prepareWorkspace(
		adapter,
		"hoily moily",
		"moily hoily is about hoily moily",
		userID,
	)
	if err != nil {
		test.Error(err)
	}

	err = makeInvitationTestCase(
		adapter,
		workspaceID,
		userID2,
		contactID,
	)
	if err != nil {
		test.Error(err)
	}

	if err = adapter.RemoveUser([]byte(entity)); err != nil {
		test.Error(err)
	}

	if err = adapter.RemoveUser([]byte(entity2)); err != nil {
		test.Error(err)
	}
}

func Test_JoinRequestRoundtrip(test *testing.T) {
	adapter, err := sqldb.NewDatabaseAdapter("keep")
	if err != nil {
		test.Error(err)
	}
	defer adapter.Close()

	userID, _, entity, err := getEssentialJoinRequestData(
		adapter,
		"supa entity",
		"dupa name",
		"supa@dupa.mupa",
	)
	if err != nil {
		test.Error(err)
	}

	_, profileID2, entity2, err := getEssentialJoinRequestData(
		adapter,
		"supa dupa entity",
		"dupa supa name",
		"supadupa@dupasupa.mupakupa",
	)
	if err != nil {
		test.Error(err)
	}

	workspaceID, err := prepareWorkspace(
		adapter,
		"julipuli",
		"kapooleyhooley",
		userID,
	)
	if err != nil {
		test.Error(err)
	}

	if err = makeJoinRequestTestCase(
		adapter,
		workspaceID,
		profileID2,
	); err != nil {
		test.Error(err)
	}

	if err = adapter.RemoveUser([]byte(entity)); err != nil {
		test.Error(err)
	}

	if err = adapter.RemoveUser([]byte(entity2)); err != nil {
		test.Error(err)
	}
}

func getEssentialInvitationData(
	adapter *sqldb.DBAdapter,
	entity string,
	username string,
	email string,
) (string, string, string, error) {
	// making user 1
	msg := &types.UserMessagePayload{
		Entity:   entity,
		Username: username,
		Email:    email,
	}

	payoad, err := json.Marshal(msg)
	if err != nil {
		return "", "", "", err
	}

	if err = adapter.AddUser(payoad); err != nil {
		return "", "", "", err
	}

	var userID sql.NullString
	var contactID sql.NullString

	if err = adapter.QueryRow("SELECT u.id, c.id FROM users u JOIN profiles p ON u.id = p.user_id JOIN contacts c ON p.id = c.profile_id;").Scan(&userID, &contactID); err != nil {
		return "", "", "", err
	}

	if !sqldb.ValidValues(userID.Valid, contactID.Valid) {
		return "", "", "", errors.New("invalid values")
	}

	if err = sqldb.UUIDsFromBin(&userID.String, &contactID.String); err != nil {
		return "", "", "", err
	}

	return userID.String, contactID.String, entity, nil
}

func prepareWorkspace(
	adapter *sqldb.DBAdapter,
	name string,
	description string,
	userID string,
) (string, error) {
	if err := adapter.InsertWorkspace(
		model.WorkspaceCredentials{
			Name:        name,
			Description: description,
			UserID:      userID,
		},
	); err != nil {
		return "", err
	}

	var workspaceID sql.NullString

	if err := adapter.QueryRow("SELECT id FROM workspaces WHERE user_id = UNHEX(REPLACE(?, '-', ''));", userID).Scan(&workspaceID); err != nil {
		return "", err
	}

	if !sqldb.ValidValues(workspaceID.Valid) {
		return "", errors.New("invalid value")
	}

	if err := sqldb.UUIDsFromBin(&workspaceID.String); err != nil {
		return "", err
	}

	return workspaceID.String, nil
}

func makeInvitationTestCase(
	adapter *sqldb.DBAdapter,
	workspaceID string,
	userID2 string,
	contactID string,
) error {
	for index := 0; index < 10; index++ {
		adapter.InsertInvitation(
			model.InvitationCredentials{
				Subject:     "halli hallo",
				Message:     "hallo halli",
				WorkspaceID: workspaceID,
				UserID:      userID2,
				ContactID:   contactID,
			},
		)
	}

	rows, err := adapter.Query("SELECT id FROM invitations;")
	if err != nil {
		return err
	}

	var invitationID sql.NullString
	var invitationIDs []any

	for rows.Next() {
		if err = rows.Scan(&invitationID); err != nil {
			return err
		}
		if !sqldb.ValidValues(invitationID.Valid) {
			return errors.New("invalid value")
		}
		if err = sqldb.UUIDsFromBin(&invitationID.String); err != nil {
			return err
		}

		invitationIDs = append(invitationIDs, invitationID.String)
	}

	log.Println(invitationIDs...)

	query := sqldb.MakeBatchQuery(
		"UPDATE invitations SET sent_at = DATE_SUB(NOW(), INTERVAL 32 DAY) WHERE id IN (%s);",
		"UNHEX(REPLACE(?, '-', ''))",
		len(invitationIDs),
	)

	log.Println(query)

	_, err = adapter.Exec(query, invitationIDs...)
	if err != nil {
		return err
	}

	ids, _, err := adapter.GetInvitations()
	if err != nil {
		return err
	}

	if len(ids) != len(invitationIDs) {
		return fmt.Errorf("expected length of %d got length of %d", len(invitationIDs), len(ids))
	}

	return adapter.RemoveSelectedInvitationData(ids)
}

func getEssentialJoinRequestData(
	adapter *sqldb.DBAdapter,
	entity string,
	username string,
	email string,
) (string, string, string, error) {
	msg := &types.UserMessagePayload{
		Entity:   entity,
		Username: username,
		Email:    email,
	}

	payoad, err := json.Marshal(msg)
	if err != nil {
		return "", "", "", err
	}

	if err = adapter.AddUser(payoad); err != nil {
		return "", "", "", err
	}

	var userID sql.NullString
	var profileID sql.NullString

	if err = adapter.QueryRow("SELECT u.id, p.id FROM users u JOIN profiles p ON u.id = p.user_id;").Scan(&userID, &profileID); err != nil {
		return "", "", "", err
	}

	if !sqldb.ValidValues(
		userID.Valid,
		profileID.Valid,
	) {
		return "", "", "", errors.New("invalid value")
	}

	if err = sqldb.UUIDsFromBin(
		&userID.String,
		&profileID.String,
	); err != nil {
		return "", "", "", err
	}

	return userID.String, profileID.String, entity, nil
}

func makeJoinRequestTestCase(
	adapter *sqldb.DBAdapter,
	workspaceID string,
	profileID2 string,
) error {
	for index := 0; index < 10; index++ {
		if err := adapter.InsertJoinRequest(model.JoinRequestCredentials{
			Subject:     "jhassajhdsjhsd",
			Message:     "kjdfkdjfdfkj",
			WorksapceID: workspaceID,
			ProfileID:   profileID2,
		},
		); err != nil {
			return err
		}
	}

	rows, err := adapter.Query("SELECT id FROM join_requests;")
	if err != nil {
		return err
	}

	var joinRequestIDs []any
	var joinRequestID sql.NullString

	for rows.Next() {
		if err := rows.Scan(&joinRequestID); err != nil {
			return err
		}

		if !sqldb.ValidValues(joinRequestID.Valid) {
			return errors.New("invalid value")
		}

		if err := sqldb.UUIDsFromBin(&joinRequestID.String); err != nil {
			return errors.New("invalid id")
		}

		joinRequestIDs = append(joinRequestIDs, joinRequestID.String)
	}

	if len(joinRequestIDs) == 0 {
		return errors.New("something went wrong with the inserts")
	}

	query := sqldb.MakeBatchQuery(
		"UPDATE join_requests SET sent_at = DATE_SUB(NOW(), INTERVAL 32 DAY) WHERE id IN (%s);",
		"UNHEX(REPLACE(?, '-', ''))",
		len(joinRequestIDs),
	)

	_, err = adapter.Exec(query, joinRequestIDs...)
	if err != nil {
		return err
	}

	ids, _, err := adapter.GetJoinRequests()
	if err != nil {
		return err
	}

	if len(ids) != len(joinRequestIDs) {
		return fmt.Errorf("expected length of %d got length of %d", len(joinRequestIDs), len(ids))
	}

	return adapter.RemoveSelectedJoinRequestData(ids)
}
