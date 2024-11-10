package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"user_manager/graph/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func makeBatchQuery(queryStatement, inputStatement string, batchSize int) string {
	batchInputs := make([]string, batchSize)

	for index := 0; index < batchSize; index++ {
		batchInputs = append(batchInputs, inputStatement)
	}

	return fmt.Sprintf(queryStatement, strings.Join(batchInputs, ","))
}

func computeExistingXPBatch(transaction *sql.Tx, names []string, numberOfEntries int) ([]*XPData, error) {
	query := makeBatchQuery(
		"SELECT e.id FROM experiences e.name IN (%s)",
		"UNHEX(REPLACE(?, \"-\", \"\"))",
		numberOfEntries,
	)
	statement, err := transaction.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	nameArgs := make([]any, len(names))

	for _, name := range names {
		nameArgs = append(nameArgs, name)
	}

	rows, err := statement.Query(nameArgs...)
	if err != nil {
		return nil, err
	}

	results := make([]*XPData, len(names))
	var id sql.NullString
	index := 0

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		} else if !id.Valid {
			return nil, err
		} else if err := uuidFromBin(&id.String); err != nil {
			return nil, err
		}

		results = append(results, &XPData{ID: id.String, Name: names[index]})
		index++
	}

	return results, nil
}

func computeStringBatchStatement(statement *sql.Stmt, ids []string) ([]string, error) {
	idArgs := make([]any, len(ids))

	for _, id := range ids {
		idArgs = append(idArgs, id)
	}

	rows, err := statement.Query(idArgs...)
	if err != nil {
		return nil, err
	}

	results := make([]string, len(ids))
	var str sql.NullString

	for rows.Next() {
		if err := rows.Scan(&str); err != nil {
			return nil, err
		} else if !str.Valid {
			return nil, errors.New("invalid id")
		} else {
			results = append(results, str.String)
		}
	}

	return results, nil
}

func computeExperienceprofileJoinBatch(transaction *sql.Tx, queryArgs []any) error {
	query := makeBatchQuery(
		"INSERT INTO experiences_profiles (experience_id, profile_id) VAlUES (%s)",
		"UNHEX(REPLACE(?, \"-\", \"\"))",
		len(queryArgs),
	)

	statement, err := transaction.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(queryArgs...); err != nil {
		return err
	}

	return nil
}

func computeNewExperienceInserts(transaction *sql.Tx, credentials []*model.NewExperienceCredentials) ([]string, error) {
	query := makeBatchQuery(
		"INSERT INTO experiences (id, name, user_id) VALUES",
		"(UNHEX(REPLACE(?, \"_\", \"\")), ?, UNHEX(REPLACE(?, \"_\", \"\")))",
		len(credentials),
	)

	statement, err := transaction.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	queryArgs := make([]any, len(credentials)*3)
	ids := make([]string, len(credentials))
	for _, experience := range credentials {
		guid := uuid.New().String()

		queryArgs = append(
			queryArgs,
			guid,
			experience.Experience,
			experience.UserID,
		)

		ids = append(ids, guid)
	}

	if _, err := statement.Exec(queryArgs...); err != nil {
		return nil, err
	}

	return ids, nil
}

func (db *Database) computeWithMixedExperiences(transaction *sql.Tx, credentials []*model.NewExperienceCredentials, xpData []*XPData) ([]*model.Experience, error) {
	count := 0
	existingExperiences := make([]*model.ExperienceCredentials, 0)
	existing := make(map[int]bool)

	for index, experience := range credentials {
		if experience.Experience == xpData[count].Name {
			existingExperiences = append(existingExperiences, &model.ExperienceCredentials{
				ExperienceID: xpData[count].ID,
				UserID:       experience.UserID,
				ProfileID:    experience.ProfileID,
				Rating:       experience.Rating,
			})
			count++

			existing[index] = true
		}
	}

	notExisting := make([]*model.NewExperienceCredentials, 0)
	for index, experience := range credentials {
		if !existing[index] {
			notExisting = append(notExisting, experience)
		}
	}

	existingResults, err := db.batchExistingExperiences(transaction, existingExperiences)
	if err != nil {
		return nil, err
	}

	newResults, err := db.computeWithNewExperiences(transaction, notExisting)
	if err != nil {
		return nil, err
	}

	results := make([]*model.Experience, len(existingResults)+len(newResults))
	results = append(results, existingResults...)
	results = append(results, newResults...)

	return results, nil
}

func (db *Database) computeWithNewExperiences(transaction *sql.Tx, credentials []*model.NewExperienceCredentials) ([]*model.Experience, error) {
	xpIDs, err := computeNewExperienceInserts(transaction, credentials)
	if err != nil {
		return nil, err
	}

	if err := computeExperienceprofileJoinBatch(
		transaction,
		getExperienceProfileJoinData(xpIDs, getProfileIDs(credentials)),
	); err != nil {
		return nil, err
	}

	ratingIDs, err := db.computeInsertRatingBatch(
		transaction,
		getaRtingsFromNewExperiences(credentials),
		credentials[0].UserID,
		xpIDs,
	)
	if err != nil {
		return nil, err
	}

	results := make([]*model.Experience, len(credentials))
	for index, experience := range credentials {
		results = append(results, &model.Experience{
			ID:         ratingIDs[index],
			Experience: experience.Experience,
			Rating:     experience.Rating,
		})
	}

	return results, nil
}

func (db *Database) computeNewExperienceBatch(transaction *sql.Tx, credentials []*model.NewExperienceCredentials) ([]*model.Experience, error) {

	names := getExperienceNames(credentials)

	xpData, err := computeExistingXPBatch(transaction, names, len(credentials))

	switch {
	case err == sql.ErrNoRows:
		return db.computeWithNewExperiences(transaction, credentials)

	case err != nil:
		return nil, err

	default:
		return db.computeWithMixedExperiences(transaction, credentials, xpData)
	}
}

func (db *Database) insertRatingBatch(transaction *sql.Tx, query string, batchData ...any) error {
	statement, err := transaction.Prepare(query)
	if err != nil {
		return err
	}

	if _, err := statement.Exec(batchData...); err != nil {
		return err
	}

	return nil
}

func (db *Database) computeInsertRatingBatch(transaction *sql.Tx, ratings []int, userID string, experienceIDs []string) ([]string, error) {
	batchData := make([]any, len(ratings)*4)
	ids := make([]string, len(ratings))

	for index, rating := range ratings {
		guid := uuid.New().String()
		batchData = append(
			batchData,
			guid,
			rating,
			experienceIDs[index],
			userID,
		)

		ids = append(ids, guid)
	}

	query := makeBatchQuery(
		"INSERT INTO ratings (id, rating, experience_id, user_id) VALUES %s",
		"(UNHEX(REPLACE(?, \"-\", \"\")), ?, UNHEX(REPLACE(?, \"-\", \"\")), UNHEX(REPLACE(?, \"-\", \"\")))",
		len(ratings),
	)

	if err := db.insertRatingBatch(transaction, query, batchData...); err != nil {
		return nil, err
	}

	return ids, nil
}

func (db *Database) batchExistingExperiences(transaction *sql.Tx, credentials []*model.ExperienceCredentials) ([]*model.Experience, error) {
	query := makeBatchQuery("SELECT e.name FROM experiences e WHERE e.id IN (%s)", "UNHEX(REPLACE(?, \"-\", \"\"))", len(credentials))

	statement, err := transaction.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	experienceIDs := getExperienceIDs(credentials)
	experienceNames, err := computeStringBatchStatement(statement, experienceIDs)
	if err != nil {
		return nil, err
	}

	ratingIDs, err := db.computeInsertRatingBatch(
		transaction,
		getRatingsFromExistingExperiences(credentials),
		credentials[0].UserID,
		getIDsFromExistingExperiences(credentials),
	)
	if err != nil {
		return nil, err
	}

	if err := computeExperienceprofileJoinBatch(
		transaction,
		getExperienceProfileJoinData(experienceIDs, getProfileIDsFromExisting(credentials)),
	); err != nil {
		return nil, err
	}

	return getExperiencesFromExisting(
		credentials,
		ratingIDs,
		experienceNames,
	), err
}
