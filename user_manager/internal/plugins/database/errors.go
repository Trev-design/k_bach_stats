package database

import "errors"

var (
	ErrNotFound                 = errors.New("entity not found")
	ErrInternal                 = errors.New("something went wrong")
	ErrConnectionLost           = errors.New("no connection after init")
	ErrForbiddenInitOpt         = errors.New("invalid init option")
	ErrForbiddenDatabaseRequest = errors.New("invalid request")
	ErrInvalidEntry             = errors.New("invalid entry in database")
	ErrInvalidBinaryID          = errors.New("invalid binary id")
)
