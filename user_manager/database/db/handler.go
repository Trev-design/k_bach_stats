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
	rows, err := db.Query(invitationInfos, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invitationInfos := make([]*model.InvitationInfo, 0)

	for rows.Next() {
		invitationInfo := new(model.InvitationInfo)
		if err := rows.Scan(
			&invitationInfo.ID,
			&invitationInfo.Info,
			&invitationInfo.UserID,
		); err != nil {
			return nil, err
		}

		if err := uuidFromBin(
			&invitationInfo.ID,
			&invitationInfo.UserID,
		); err != nil {
			return nil, err
		}

		invitationInfos = append(invitationInfos, invitationInfo)
	}

	return invitationInfos, nil
}

func (db *Database) GetJoinRequestInfosFromDB(workspaceID string) ([]*model.JoinRequestInfo, error) {
	rows, err := db.Query(joinRequestInfos, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]*model.JoinRequestInfo, 0)

	for rows.Next() {
		request := new(model.JoinRequestInfo)

		if err := rows.Scan(
			&request.ID,
			&request.Info,
			&request.JoinRequestID,
		); err != nil {
			return nil, err
		}

		if err := uuidFromBin(
			&request.ID,
			&request.JoinRequestID,
		); err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (db *Database) GetWorkspaceFromDB(workspaceID string) (*model.CompleteWorkspace, error) {
	rows, err := db.Query(completeWorkspace, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workspace := &model.CompleteWorkspace{
		Contacts: make([]*model.Contact, 0),
		News:     make([]*model.InvitationInfo, 0),
	}

	var contactID sql.NullString
	var contactName sql.NullString
	var contactEmail sql.NullString
	var contactImagefilePath sql.NullString

	var invitationInfoID sql.NullString
	var invitationInfoInfo sql.NullString
	var invitationInfoInvitationID sql.NullString

	for rows.Next() {
		if err := rows.Scan(
			&workspace.ID,
			&workspace.Name,
			&workspace.Description,
			&contactID,
			&contactName,
			&contactEmail,
			&contactImagefilePath,
			&invitationInfoID,
			&invitationInfoInfo,
			&invitationInfoInvitationID,
		); err != nil {
			return nil, err
		}

		if validvalues(
			contactID,
			contactName,
			contactEmail,
			contactImagefilePath,
		) {
			contact := &model.Contact{
				ID:            contactID.String,
				Name:          contactName.String,
				Email:         contactEmail.String,
				ImageFilePath: contactImagefilePath.String,
			}

			if err := uuidFromBin(&contact.ID); err != nil {
				return nil, err
			}

			workspace.Contacts = append(workspace.Contacts, contact)
		}

		if validvalues(
			invitationInfoID,
			invitationInfoInfo,
			invitationInfoInvitationID,
		) {
			invitationInfo := &model.InvitationInfo{
				ID:     invitationInfoID.String,
				Info:   invitationInfoInfo.String,
				UserID: invitationInfoInvitationID.String,
			}

			if err := uuidFromBin(
				&invitationInfo.ID,
				&invitationInfo.UserID,
			); err != nil {
				return nil, err
			}

			workspace.News = append(workspace.News, invitationInfo)
		}
	}

	return workspace, nil
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

func (db *Database) UpdateBio(credentials model.BioCredentials) error {
	if _, err := db.Exec(
		updateBio,
		credentials.Input,
		credentials.UserID,
	); err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateName(credentials model.ChangeNameCredentials) error {
	if _, err := db.Exec(
		updateName,
		credentials.Newname,
		credentials.UserID,
	); err != nil {
		return err
	}

	return nil
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

func validvalues(strings ...sql.NullString) bool {
	for _, val := range strings {
		if !val.Valid {
			return false
		}
	}

	return true
}
