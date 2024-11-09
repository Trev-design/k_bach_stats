package database

import "user_manager/graph/model"

type UserHandler interface {
	AddUser(payload []byte) error
	RemoveUser(payload []byte) error
}

type SessionHandler interface {
	CheckSession(token string) error
	InitialAuth(token string) (string, error)
}

type StoreHandler interface {
	InitialCredentials(entity string) (string, error)
	GetUserFromDB(entity string) (*model.User, error)
	GetInvitationInfosFromDB(userID string) ([]*model.InvitationInfo, error)
	GetJoinRequestInfosFromDB(workspaceID string) ([]*model.JoinRequestInfo, error)
	GetWorkspaceFromDB(workspaceID string) (*model.CompleteWorkspace, error)
	CreateNewWorkspace(credentials model.WorkspaceCredentials) error
	PushInvitation(credentials model.InvitationCredentials) error
	PushJoinRequest(credentials model.JoinRequestCredentials) error
	UpdateBio(credentials model.BioCredentials) error
	UpdateName(credentials model.ChangeNameCredentials) error
	NewExperience(credentials *model.NewExperienceCredentials) (*model.Experience, error)
	AddExperience(credentials *model.ExperienceCredentials) (*model.Experience, error)
	AddExperienceBatch(credentials *model.ExperienceBatchCredentials) ([]*model.Experience, error)
}
