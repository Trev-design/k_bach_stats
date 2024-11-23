package core

import "user_manager/graph/model"

type ApiService interface {
	GetUserID(entity string) (*model.UserEntity, error)
	GetUserByID(userID string) (*model.User, error)
	Wait()
	Close() error
}
