package core

import "user_manager/graph/model"

type Database interface {
	UserID(entity string) (*model.UserEntity, error)
	UserByID(userID string) (*model.User, error)
	Cancel()
	Closed() bool
	CloseDB() error
	Wait()
}
