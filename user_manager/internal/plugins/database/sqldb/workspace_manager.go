package sqldb

import (
	"user_manager/graph/model"

	"github.com/google/uuid"
)

func (dba *DBAdapter) InsertWorkspace(input model.WorkspaceCredentials) error {
	guid := uuid.New().String()
	_, err := dba.Exec(
		insertWorkspaceQuery,
		guid,
		input.Name,
		input.Description,
		input.UserID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (dba *DBAdapter) RemoveWorkspace(workspaceID string) error {
	_, err := dba.Exec("DELETE FROM workspaces WHERE id = ?;", workspaceID)
	if err != nil {
		return err
	}

	return nil
}
