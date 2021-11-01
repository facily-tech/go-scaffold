package models

import "errors"

var (
	// ErrNew represents an error when creating a new Quote
	ErrNew = errors.New("unable to create new quote")
)
