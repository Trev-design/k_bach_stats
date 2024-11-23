package application

import (
	"sync"
	"user_manager/graph/model"
	"user_manager/internal/core"
)

type APIServiceAdapter struct {
	db core.Database
	wg *sync.WaitGroup
}

func NewAPIServiceAdapter(db core.Database) *APIServiceAdapter {
	return &APIServiceAdapter{db: db}
}

func (apia *APIServiceAdapter) GetUserID(entity string) (*model.UserEntity, error) {
	return apia.db.UserID(entity)
}

func (apia *APIServiceAdapter) GetUserByID(userID string) (*model.User, error) {
	return apia.db.UserByID(userID)
}

func (apia *APIServiceAdapter) Wait() {
	apia.wg.Wait()
}

func (apia *APIServiceAdapter) Close() error {
	apia.db.Cancel()
	apia.db.Wait()

	return apia.db.CloseDB()
}
