package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

type XPData struct {
	ID   string
	Name string
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
		Experiences: make([]*model.Experience, 0),
	}

	var experienceName sql.NullString
	var rating sql.NullInt32
	var ratingID sql.NullString

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
			&experienceName,
			&rating,
			&ratingID,
		); err != nil {
			return nil, err
		}

		if experienceName.Valid && rating.Valid && ratingID.Valid {
			experience := &model.Experience{
				Experience: experienceName.String,
				Rating:     int(rating.Int32),
				ID:         ratingID.String,
			}
			user.Experiences = append(user.Experiences, experience)
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
			&invitationInfo.InvitationID,
		); err != nil {
			return nil, err
		}

		if err := uuidFromBin(
			&invitationInfo.ID,
			&invitationInfo.InvitationID,
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
				ID:           invitationInfoID.String,
				Info:         invitationInfoInfo.String,
				InvitationID: invitationInfoInvitationID.String,
			}

			if err := uuidFromBin(
				&invitationInfo.ID,
				&invitationInfo.InvitationID,
			); err != nil {
				return nil, err
			}

			workspace.News = append(workspace.News, invitationInfo)
		}
	}

	return workspace, nil
}

func (db *Database) CreateNewWorkspace(credentials model.WorkspaceCredentials) error {
	guid := uuid.New().String()

	if _, err := db.Exec(
		createWorkspaceQuery,
		guid,
		credentials.Name,
		credentials.Description,
		credentials.UserID,
	); err != nil {
		return err
	}

	return nil
}

func (db *Database) PushInvitation(credentials model.InvitationCredentials) error {
	guid := uuid.New().String()
	if _, err := db.Exec(
		createInvitationQuery,
		guid,
		credentials.Info,
		credentials.InvitorID,
		credentials.ReceiverID,
		credentials.WorkspaceID,
	); err != nil {
		return err
	}

	return nil
}

func (db *Database) PushJoinRequest(credentials model.JoinRequestCredentials) error {
	guid := uuid.New().String()

	if _, err := db.Exec(
		createRequestQuery,
		guid,
		credentials.Info,
		credentials.Reason,
		credentials.WorkspaceID,
		credentials.RequestID,
	); err != nil {
		return err
	}

	return nil
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

func (db *Database) AddExperience(credentials *model.ExperienceCredentials) (*model.Experience, error) {
	transaction, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer handleTransaction(transaction, err)

	result, err := db.insertExistingExperience(transaction, credentials)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *Database) NewExperience(credentials *model.NewExperienceCredentials) (*model.Experience, error) {
	transaction, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer handleTransaction(transaction, err)

	result, err := db.computeNewExperience(transaction, credentials)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *Database) AddExperienceBatch(credentials *model.ExperienceBatchCredentials) ([]*model.Experience, error) {
	transaction, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer handleTransaction(transaction, err)

	results := []*model.Experience{}

	newResults, err := db.computeNewExperienceBatch(transaction, credentials.New)
	if err != nil {
		return nil, err
	}

	results = append(results, newResults...)

	existingResults, err := db.batchExistingExperiences(transaction, credentials.Existing)
	if err != nil {
		return nil, err
	}

	results = append(results, existingResults...)

	return results, nil
}

func (db *Database) insertRating(transaction *sql.Tx, rating int, userID, experienceID string) (string, error) {
	statement, err := transaction.Prepare(addExperienceQuery)
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

func (db *Database) insertExistingExperience(transaction *sql.Tx, credentials *model.ExperienceCredentials) (*model.Experience, error) {
	statement, err := transaction.Prepare(experienceByIDQuery)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var name sql.NullString

	if statement.QueryRow(credentials.ExperienceID).Scan(name); err != nil {
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

	return &model.Experience{
		ID:         ratingID,
		Experience: name.String,
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

		return &model.Experience{
			ID:         ratingID,
			Experience: credentials.Experience,
			Rating:     credentials.Rating,
		}, nil
	}

	return db.insertNewExperience(transaction, credentials)
}

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

	return &model.Experience{
		ID:         ratingID,
		Experience: credentials.Experience,
		Rating:     credentials.Rating,
	}, nil
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

	return getExperiencesFromExisting(
		credentials,
		ratingIDs,
		experienceNames,
	), err
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

func (db *Database) computeNewExperienceBatch(transaction *sql.Tx, credentials []*model.NewExperienceCredentials) ([]*model.Experience, error) {
	query := makeBatchQuery(
		"SELECT e.id FROM experiences e.name IN (%s)",
		"UNHEX(REPLACE(?, \"-\", \"\"))",
		len(credentials),
	)
	statement, err := transaction.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	names := getExperienceNames(credentials)

	xpData, err := computeExistingXPBatch(statement, names)

	switch {
	case err == sql.ErrNoRows:
		return db.computeWithNewExperiences(transaction, credentials)

	case err != nil:
		return nil, err

	default:
		return db.computeWithMixedExperiences(transaction, credentials, xpData)
	}
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

func (db *Database) computeWithMixedExperiences(transaction *sql.Tx, credentials []*model.NewExperienceCredentials, xpData []*XPData) ([]*model.Experience, error) {

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

func getExperienceProfileJoinData(experiencIDs []string, profileIDs []string) []any {
	joinData := make([]any, len(experiencIDs))

	for index, experienceID := range experiencIDs {
		joinData = append(
			joinData,
			experienceID,
			profileIDs[index],
		)
	}

	return joinData
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

func computeExistingXPBatch(statement *sql.Stmt, names []string) ([]*XPData, error) {
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

func makeBatchQuery(queryStatement, inputStatement string, batchSize int) string {
	batchInputs := make([]string, batchSize)

	for index := 0; index < batchSize; index++ {
		batchInputs = append(batchInputs, inputStatement)
	}

	return fmt.Sprintf(queryStatement, strings.Join(batchInputs, ","))
}

func getExperiencesFromExisting(existing []*model.ExperienceCredentials, ratingIDs []string, names []string) []*model.Experience {
	results := make([]*model.Experience, len(existing))

	for index, experience := range existing {
		results = append(
			results,
			&model.Experience{
				ID:         ratingIDs[index],
				Experience: names[index],
				Rating:     experience.Rating,
			},
		)
	}

	return results
}

func getExperienceIDs(experiences []*model.ExperienceCredentials) []string {
	ids := make([]string, len(experiences))

	for _, experience := range experiences {
		ids = append(ids, experience.ExperienceID)
	}

	return ids
}

func getExperienceNames(experiences []*model.NewExperienceCredentials) []string {
	names := make([]string, len(experiences))

	for _, experience := range experiences {
		names = append(names, experience.Experience)
	}

	return names
}

func getProfileIDs(experiences []*model.NewExperienceCredentials) []string {
	ids := make([]string, len(experiences))

	for _, experience := range experiences {
		ids = append(ids, experience.ProfileID)
	}

	return ids
}

func getRatingsFromExistingExperiences(experiences []*model.ExperienceCredentials) []int {
	ratings := make([]int, len(experiences))
	for _, experience := range experiences {
		ratings = append(ratings, experience.Rating)
	}

	return ratings
}

func getaRtingsFromNewExperiences(experiences []*model.NewExperienceCredentials) []int {
	ratings := make([]int, len(experiences))
	for _, experience := range experiences {
		ratings = append(ratings, experience.Rating)
	}

	return ratings
}

func getIDsFromExistingExperiences(experiences []*model.ExperienceCredentials) []string {
	ids := make([]string, len(experiences))
	for _, experience := range experiences {
		ids = append(ids, experience.ExperienceID)
	}

	return ids
}

func getNewIDsForExperiences(numberOfIDs int) []string {
	ids := make([]string, numberOfIDs)
	for index := 0; index < numberOfIDs; index++ {
		ids = append(ids, uuid.New().String())
	}

	return ids
}

func handleTransaction(transaction *sql.Tx, err error) {
	if p := recover(); p != nil {
		transaction.Rollback()
	} else if err != nil {
		transaction.Rollback()
	} else if commitError := transaction.Commit(); commitError != nil {
		transaction.Rollback()
	}
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
