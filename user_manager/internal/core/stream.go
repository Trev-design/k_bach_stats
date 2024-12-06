package core

import "user_manager/types"

type Stream interface {
	SendStream(payload *types.StreamPayload) error
	Close() error
}
