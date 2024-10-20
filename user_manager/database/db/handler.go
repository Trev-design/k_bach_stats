package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"user_manager/graph/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Database struct {
	*sql.DB
}

type UserPayload struct {
	Entity   string `json:"entity"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (db *Database) AddUser(payload []byte) error {
	user := new(UserPayload)
	if err := json.Unmarshal(payload, user); err != nil {
		return err
	}

	kinds := []string{"user", "profile", "contact"}
	foreignKey := ""

	for _, kind := range kinds {
		guid, err := db.insertUserCredentials(kind, foreignKey, user)
		if err != nil {
			return err
		}

		foreignKey = guid
	}

	return nil
}

func (db *Database) RemoveUser(payload []byte) error {
	deletable := &struct {
		Entity string `json:"entity"`
	}{}
	if err := json.Unmarshal(payload, deletable); err != nil {
		fmt.Printf("could not parse from json: %v", err)
		return err
	}
	return db.removeItem(removeUser, deletable.Entity)
}

func (db *Database) InitialCredentials(entity string) (string, error) {
	var id sql.NullString
	if err := db.QueryRow(userCredentials, entity).Scan(&id); err != nil {
		return "", err
	}

	if !id.Valid {
		return "", errors.New("no user found")
	}

	return id.String, nil
}

func (db *Database) GetUserFromDB(entity string) (*model.User, error) {
	rows, err := db.Query(selectUser, entity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &model.User{
		Profile: &model.Profile{
			Contact: &model.Contact{},
		},
		Workspaces: make([]*model.Workspace, 0),
	}

	var workspace_id sql.NullString
	var workspace_name sql.NullString

	for rows.Next() {
		if err := rows.Scan(
			&user.ID,
			&user.Entity,
			&user.Requests,
			&user.Profile.ID,
			&user.Profile.Bio,
			&user.Profile.Contact.ID,
			&user.Profile.Contact.Name,
			&user.Profile.Contact.Email,
			&user.Profile.Contact.ImageFilePath,
			&workspace_id,
			&workspace_name,
		); err != nil {
			return nil, err
		}

		if workspace_id.Valid && workspace_name.Valid {
			workspace := &model.Workspace{ID: workspace_id.String, Name: workspace_name.String}
			if err := uuidFromBin(&workspace.ID); err != nil {
				return nil, err
			}
			user.Workspaces = append(user.Workspaces, workspace)
		}
	}

	if err := uuidFromBin(&user.ID, &user.Profile.ID, &user.Profile.Contact.ID); err != nil {
		return nil, err
	}

	return user, nil
}

func (db *Database) GetInvitationInfosFromDB(userID string) ([]*model.InvitationInfo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (db *Database) GetJoinRequestInfosFromDB(workspaceID string) ([]*model.JoinRequestInfo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (db *Database) GetWorkspaceFromDB(workspaceID string) (*model.CompleteWorkspace, error) {
	panic(fmt.Errorf("not implemented"))
}

func (db *Database) CreateNewWorkspace(credentials model.WorkspaceCredentials) error {
	panic(fmt.Errorf("not implemented"))
}

func (db *Database) PushInvitation(credentials model.InvitationCredentials) error {
	panic(fmt.Errorf("not implemented"))
}

func (db *Database) PushJoinRequest(credentials model.JoinRequestCredentials) error {
	panic(fmt.Errorf("not implemented"))
}

func uuidFromBin(binaryIDs ...*string) error {
	for _, id := range binaryIDs {
		guid, err := uuid.FromBytes([]byte(*id))
		if err != nil {
			return err
		}

		*id = guid.String()
	}

	return nil
}

func (db *Database) insertUserCredentials(kind, foreignKey string, user *UserPayload) (string, error) {
	guid := uuid.New().String()

	switch kind {
	case "user":
		return guid, db.insertItem(insertNewUser, guid, user.Entity)
	case "profile":
		return guid, db.insertItem(insertNewProfile, guid, "", foreignKey)
	case "contact":
		return guid, db.insertItem(insertNewContact, guid, user.Username, user.Email, "", foreignKey)
	default:
		return "", errors.New("invalid insert option")
	}
}

func (db *Database) insertItem(query string, args ...any) error {
	if _, err := db.Exec(query, args...); err != nil {
		return err
	}
	return nil
}

func (db *Database) removeItem(query string, args ...any) error {
	if _, err := db.Exec(query, args...); err != nil {
		return err
	}
	return nil
}
