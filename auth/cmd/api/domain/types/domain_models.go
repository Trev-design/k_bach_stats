package types

import "github.com/google/uuid"

type NewAccountDM struct {
	Name         string
	Email        string
	PasswordHash string
}

type UpdateStateDM struct {
	ID            uuid.UUID
	VerifiedState bool
	AboType       string
}

type AccountDM struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Isverified   bool
	AboType      string
}
