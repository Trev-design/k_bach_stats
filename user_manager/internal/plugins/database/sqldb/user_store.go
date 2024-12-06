package sqldb

import (
	"database/sql"
	"log"
	"user_manager/graph/model"
	"user_manager/internal/plugins/database"

	_ "github.com/go-sql-driver/mysql"
)

func (dba *DBAdapter) UserID(entity string) (*model.UserEntity, error) {
	var id sql.NullString

	if err := dba.QueryRow(userIDQuery, entity).Scan(&id); err != nil {
		return nil, err
	}

	if !id.Valid {
		return nil, database.ErrInvalidEntry
	}

	log.Println(id.String)

	userID, err := uuidFromBin(id.String)
	if err != nil {
		return nil, err
	}

	return &model.UserEntity{User: userID}, nil
}

func (dba *DBAdapter) UserByID(id string) (*model.User, error) {
	var profileID sql.NullString
	var bio sql.NullString
	var contactID sql.NullString
	var name sql.NullString
	var email sql.NullString
	var hasImage sql.NullString

	if err := dba.QueryRow(userByIDQuery, id).Scan(
		&profileID,
		&bio,
		&contactID,
		&name,
		&email,
		&hasImage,
	); err != nil {
		log.Println(err)
		return nil, database.ErrInvalidEntry
	}

	if !validValues(
		profileID.Valid,
		bio.Valid,
		contactID.Valid,
		name.Valid,
		email.Valid,
		hasImage.Valid,
	) {
		return nil, database.ErrInvalidEntry
	}

	if err := uuidsFrombin(
		&profileID.String,
		&contactID.String,
	); err != nil {
		log.Println(err)
		return nil, database.ErrInvalidBinaryID
	}

	return &model.User{
		ID: id,
		Profile: &model.Profile{
			ID:  profileID.String,
			Bio: bio.String,
			Contact: &model.Contact{
				ID:       contactID.String,
				Name:     name.String,
				Email:    email.String,
				HasImage: hasImage.String,
			},
		},
	}, nil
}

func (dba *DBAdapter) InsertInvitation(input model.InvitationCredentials) error {
	return dba.isnertInvitation(input)
}

func (dba *DBAdapter) InsertJoinRequest(input model.JoinRequestCredentials) error {
	return dba.insertJoinRequest(input)
}
