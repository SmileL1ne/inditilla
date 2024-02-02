package entity

import "errors"

var (
	ErrNoRecord           = errors.New("entity: no matching row found")
	ErrDuplicateEmail     = errors.New("entity: duplicate email")
	ErrInvalidCredentials = errors.New("entity: invalid credentials")
)

type ErrorResponse struct {
	Message string
	Details map[string]string
}
