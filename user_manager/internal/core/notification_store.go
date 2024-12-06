package core

import "user_manager/types"

type NotificationStore interface {
	GetInvitations() ([]any, []*types.InvitationArchivePayload, error)
	GetJoinRequests() ([]any, []*types.JoinRequestArchivePayload, error)
	RemoveSelectedInvitationData(ids []any) error
	RemoveSelectedJoinRequestData(ids []any) error
}
