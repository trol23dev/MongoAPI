package domain

import "errors"

var (
	ErrNotFound       = errors.New("document not found")
	ErrInvalidID      = errors.New("invalid document id")
	ErrEmptyBody      = errors.New("request body cannot be empty")
	ErrDatabaseEmpty  = errors.New("database name is required")
	ErrCollectionEmpty = errors.New("collection name is required")
)
