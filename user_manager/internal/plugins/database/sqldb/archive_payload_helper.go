package sqldb

import (
	"database/sql"
	"errors"
	"user_manager/types"

	_ "github.com/go-sql-driver/mysql"
)

type rawInvitation struct {
	invitationID string
	subject      string
	message      string
	userID       string
	workspaceID  string
	contactID    string
}

type rawJoinRequest struct {
	subject       string
	message       string
	profileID     string
	workspaceID   string
	joinRequestID string
}

func getRawInvitationData(transaction *sql.Tx) ([]*rawInvitation, error) {
	statement, err := transaction.Prepare(selectInvitationsQuery)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query()
	if err != nil {
		return nil, err
	}

	var invitationID sql.NullString
	var subject sql.NullString
	var message sql.NullString
	var workspaceID sql.NullString
	var userID sql.NullString
	var contactID sql.NullString

	var invitations []*rawInvitation

	for rows.Next() {
		if err := rows.Scan(
			&invitationID,
			&subject,
			&message,
			&workspaceID,
			&userID,
			&contactID,
		); err != nil {
			return nil, err
		}

		if !validValues(
			invitationID.Valid,
			subject.Valid,
			message.Valid,
			workspaceID.Valid,
			userID.Valid,
			contactID.Valid,
		) {
			return nil, errors.New("invalid values")
		}

		if err := uuidsFrombin(
			&invitationID.String,
			&workspaceID.String,
			&userID.String,
			&contactID.String,
		); err != nil {
			return nil, err
		}

		invitations = append(
			invitations,
			&rawInvitation{
				subject:      subject.String,
				message:      message.String,
				contactID:    contactID.String,
				workspaceID:  workspaceID.String,
				userID:       userID.String,
				invitationID: invitationID.String,
			},
		)
	}

	return invitations, nil
}

func getRawJoinRequestdata(transaction *sql.Tx) ([]*rawJoinRequest, error) {
	statement, err := transaction.Prepare(selectJoinRequestsQuery)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query()
	if err != nil {
		return nil, err
	}

	var requestID sql.NullString
	var subject sql.NullString
	var message sql.NullString
	var profileId sql.NullString
	var workspaceID sql.NullString

	var joinRequests []*rawJoinRequest

	for rows.Next() {
		if err = rows.Scan(
			&requestID,
			&subject,
			&message,
			&profileId,
			&workspaceID,
		); err != nil {
			return nil, err
		}

		if !validValues(
			requestID.Valid,
			subject.Valid,
			message.Valid,
			profileId.Valid,
			workspaceID.Valid,
		) {
			return nil, errors.New("invalid values")
		}

		if err = uuidsFrombin(
			&requestID.String,
			&profileId.String,
			&workspaceID.String,
		); err != nil {
			return nil, err
		}

		joinRequests = append(
			joinRequests,
			&rawJoinRequest{
				subject:       subject.String,
				message:       message.String,
				profileID:     profileId.String,
				workspaceID:   workspaceID.String,
				joinRequestID: requestID.String,
			},
		)
	}

	return joinRequests, nil
}

func getInvitationCredentials(transaction *sql.Tx, payload []*rawInvitation) ([]any, []*types.InvitationArchivePayload, error) {
	invitorIDs := make([]any, len(payload))
	userIDs := make([]any, len(payload))
	workspaceIDs := make([]any, len(payload))
	ids := make([]any, len(payload))

	for index, item := range payload {
		invitorIDs[index] = item.contactID
		userIDs[index] = item.userID
		workspaceIDs[index] = item.workspaceID
		ids[index] = item.invitationID
	}

	invitors, err := getItems(
		"SELECT c.name FROM users u JOIN profiles p ON u.id = p.user_id JOIN contacts c ON p.id = c.profile_id WHERE u.id IN (%s)",
		"UNHEX(REPLACE(?, '-', ''))",
		transaction,
		invitorIDs,
	)
	if err != nil {
		return nil, nil, err
	}

	names, err := getItems(
		"SELECT c.name FROM contacts c WHERE c.id IN (%s)",
		"UNHEX(REPLACE(?, '-', ''))",
		transaction,
		userIDs,
	)
	if err != nil {
		return nil, nil, err
	}

	workspaceNames, err := getItems(
		"SELECT w.name FROM workspaces w WHERE w.id (%s)",
		"UNHEX(REPLACE(?, '-', ''))",
		transaction,
		workspaceIDs,
	)
	if err != nil {
		return nil, nil, err
	}

	items := make([]*types.InvitationArchivePayload, len(payload))

	for index, item := range payload {
		items[index] = &types.InvitationArchivePayload{
			Message:      item.message,
			Subject:      item.subject,
			ID:           item.userID,
			InvitorID:    item.contactID,
			Invitor:      invitors[index],
			InvitedGuest: names[index],
			Workspace:    workspaceNames[index],
		}
	}

	return ids, items, nil
}

func getJoinRequestCredentials(transaction *sql.Tx, payload []*rawJoinRequest) ([]any, []*types.JoinRequestArchivePayload, error) {
	ids := make([]any, len(payload))
	profileIDs := make([]any, len(payload))
	workspaceIDs := make([]any, len(payload))

	for index, item := range payload {
		ids[index] = item.joinRequestID
		profileIDs[index] = item.profileID
		workspaceIDs[index] = item.workspaceID
	}

	requesterNames, err := getItems(
		"SELECT c.name FROM profiles p JOIN contacts c ON p.id = c.profile_id WHERE p.id IN (%s)",
		"UNHEX(REPLACE(?, '-', ''))",
		transaction,
		profileIDs,
	)
	if err != nil {
		return nil, nil, err
	}

	workspaceNames, err := getItems(
		"SELECT w.name FROM workspaces w WHERE w.id (%s)",
		"UNHEX(REPLACE(?, '-', ''))",
		transaction,
		workspaceIDs,
	)
	if err != nil {
		return nil, nil, err
	}

	items := make([]*types.JoinRequestArchivePayload, len(payload))

	for index, item := range payload {
		items[index] = &types.JoinRequestArchivePayload{
			Message:     item.message,
			Subject:     item.subject,
			ID:          item.workspaceID,
			RequesterID: item.profileID,
			Requester:   requesterNames[index],
			Workspace:   workspaceNames[index],
		}
	}

	return ids, items, nil
}

func getItems(queryStatement, placeholder string, transaction *sql.Tx, ids []any) ([]string, error) {
	items := make([]string, len(ids))
	batchSize := 100

	for index := 0; index+batchSize < len(ids); index += batchSize {
		end := index + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[index:end]

		pack, err := getNames(queryStatement, placeholder, transaction, batch)
		if err != nil {
			return nil, err
		}

		items = append(items, pack...)
	}

	return items, nil
}

func deleteItems(batchQueryStatement, placeholder string, transaction *sql.Tx, ids []any) error {
	batchSize := 100

	for index := 0; index+batchSize < len(ids); index += batchSize {
		end := index + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[index:end]

		if err := delete(
			batchQueryStatement,
			placeholder,
			transaction,
			batch,
		); err != nil {
			return err
		}
	}

	return nil
}

func delete(queryStatement, placeholder string, transaction *sql.Tx, ids []any) error {
	query := makeBatchQuery(
		queryStatement,
		placeholder,
		len(ids),
	)

	statement, err := transaction.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(ids...); err != nil {
		return err
	}

	return nil
}

func getNames(queryStatement, placeholder string, transaction *sql.Tx, contactIDs []any) ([]string, error) {
	query := makeBatchQuery(
		queryStatement,
		placeholder,
		len(contactIDs),
	)

	statement, err := transaction.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(contactIDs...)
	if err != nil {
		return nil, err
	}

	var name sql.NullString
	count := 0
	names := make([]string, len(contactIDs))

	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		if !name.Valid {
			return nil, errors.New("invalid value")
		}

		names[count] = name.String
		count++
	}

	return names, nil
}
