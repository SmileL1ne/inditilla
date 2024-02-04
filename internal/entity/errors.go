package entity

import "errors"

var (
	ErrNoRecord           = errors.New("entity: no matching row found")
	ErrDuplicateEmail     = errors.New("entity: duplicate email")
	ErrInvalidCredentials = errors.New("entity: invalid credentials")
	ErrInvalidAccessToken = errors.New("invalid auth token")
)

type ErrorResponse struct {
	ResponseStatus string `json:"responseStatus"`
	Code           int    `json:"code"`
	Message        string `json:"message"`
	Location       string `json:"location,omitempty"`
}
