package db

import (
	"database/sql"
	"errors"
	"user_manager/graph/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func (db *Database) insertNewExperience(transaction *sql.Tx, credentials *model.NewExperienceCredentials) (*model.Experience, error) {
	statement, err := transaction.Prepare(addExperienceQuery)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	guid := uuid.New().String()

	ratingID, err := db.insertRating(
		transaction,
		credentials.Rating,
		credentials.UserID,
		guid,
	)
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	if err := db.experienceProfileJoin(
		transaction,
		guid,
		credentials.ProfileID,
	); err != nil {
		return nil, err
	}

	return &model.Experience{
		ID:         ratingID,
		Experience: credentials.Experience,
		Rating:     credentials.Rating,
	}, nil
}

func (db *Database) computeNewExperience(transaction *sql.Tx, credentials *model.NewExperienceCredentials) (*model.Experience, error) {
	statement, err := transaction.Prepare(experienceByNameQuery)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var id sql.NullString
	if err := statement.QueryRow(credentials.Experience).Scan(id); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if id.Valid {
		if err := uuidFromBin(&id.String); err != nil {
			return nil, err
		}

		ratingID, err := db.insertRating(
			transaction,
			credentials.Rating,
			credentials.UserID,
			id.String,
		)
		if err != nil {
			return nil, errors.New("something went wrong")
		}

		if err := db.experienceProfileJoin(
			transaction,
			id.String,
			credentials.ProfileID,
		); err != nil {
			return nil, err
		}

		return &model.Experience{
			ID:         ratingID,
			Experience: credentials.Experience,
			Rating:     credentials.Rating,
		}, nil
	}

	return db.insertNewExperience(transaction, credentials)
}

func (db *Database) insertExistingExperience(transaction *sql.Tx, credentials *model.ExperienceCredentials) (*model.Experience, error) {
	statement, err := transaction.Prepare(experienceByIDQuery)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var name sql.NullString

	if err := statement.QueryRow(credentials.ExperienceID).Scan(name); err != nil {
		return nil, err
	}

	if !name.Valid {
		return nil, errors.New("not existing experience were used")
	}

	ratingID, err := db.insertRating(
		transaction,
		credentials.Rating,
		credentials.UserID,
		credentials.ExperienceID,
	)
	if err != nil {
		return nil, err
	}

	if err := db.experienceProfileJoin(
		transaction,
		credentials.ExperienceID,
		credentials.ProfileID); err != nil {
		return nil, err
	}

	return &model.Experience{
		ID:         ratingID,
		Experience: name.String,
		Rating:     credentials.Rating,
	}, nil
}

func (db *Database) insertRating(transaction *sql.Tx, rating int, userID, experienceID string) (string, error) {
	statement, err := transaction.Prepare(addRatingQuery)
	if err != nil {
		return "", err
	}
	defer statement.Close()

	guid := uuid.New().String()

	_, err = statement.Exec(guid, rating, experienceID, userID)
	if err != nil {
		return "", err
	}

	return guid, nil
}

func (db *Database) experienceProfileJoin(transaction *sql.Tx, experienceID, profileID string) error {
	statement, err := transaction.Prepare(addProfileExperienceJoinItemQuery)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(experienceID, profileID); err != nil {
		return err
	}

	return nil
}
