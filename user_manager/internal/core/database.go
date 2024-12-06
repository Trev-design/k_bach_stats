package core

import "user_manager/graph/model"

type Database interface {
	UserID(entity string) (*model.UserEntity, error)
	UserByID(userID string) (*model.User, error)
	InsertInvitation(input model.InvitationCredentials) error
	InsertJoinRequest(input model.JoinRequestCredentials) error
	//InsertWorkspace()
	Cancel()
	Closed() bool
	CloseDB() error
	Wait()
}
