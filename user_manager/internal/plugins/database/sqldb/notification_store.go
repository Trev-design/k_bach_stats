package sqldb

import "user_manager/types"

func (dba *DBAdapter) GetInvitations() ([]any, []*types.InvitationArchivePayload, error) {
	transaction, err := dba.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer handleTransaction(transaction, err)

	rawInvitations, err := getRawInvitationData(transaction)
	if err != nil {
		return nil, nil, err
	}

	ids, items, err := getInvitationCredentials(transaction, rawInvitations)
	if err != nil {
		return nil, nil, err
	}

	return ids, items, nil
}

func (dba *DBAdapter) GetJoinRequests() ([]any, []*types.JoinRequestArchivePayload, error) {
	transaction, err := dba.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer handleTransaction(transaction, err)

	rawRequests, err := getRawJoinRequestdata(transaction)
	if err != nil {
		return nil, nil, err
	}

	ids, items, err := getJoinRequestCredentials(transaction, rawRequests)
	if err != nil {
		return nil, nil, err
	}

	return ids, items, nil
}

func (dba *DBAdapter) RemoveSelectedInvitationData(ids []any) error {
	transaction, err := dba.Begin()
	if err != nil {
		return err
	}
	defer handleTransaction(transaction, err)

	err = deleteItems(
		"DELETE FROM invitations i WHERE i.id IN (%s)",
		"UNHEX(REPLACE(?, '-', ''))",
		transaction,
		ids,
	)
	if err != nil {
		return err
	}

	return nil
}

func (dba *DBAdapter) RemoveSelectedJoinRequestData(ids []any) error {
	transaction, err := dba.Begin()
	if err != nil {
		return err
	}
	defer handleTransaction(transaction, err)

	err = deleteItems(
		"DELETE FROM join_requests j WHERE j.id IN (%s)",
		"UNHEX(REPLACE(?, '-', ''))",
		transaction,
		ids,
	)
	if err != nil {
		return err
	}

	return nil
}
