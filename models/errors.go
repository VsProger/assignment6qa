package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")

	ErrDuplicateName = errors.New("models: duplicate name")

	UnknownCategory = errors.New("models: category doesnt exist")
)
