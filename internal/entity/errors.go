package entity

import "errors"

var (
	ErrNoRecord           = errors.New("entity: no matching row found")
	ErrDuplicateEmail     = errors.New("entity: duplicate email")
	ErrInvalidCredentials = errors.New("entity: invalid credentials")
	ErrInvalidInputData   = errors.New("entity: invalid form fill")
	ErrInvalidUserId      = errors.New("entity: invalid user id")
	ErrInvalidAccessToken = errors.New("entity: invalid auth token")
	ErrEditConflict       = errors.New("entity: edit conflict")
)

type ErrorResponse struct {
	ResponseStatus string            `json:"responseStatus"`
	Code           int               `json:"code"`
	Message        string            `json:"message"`
	Location       string            `json:"location,omitempty"`
	Validations    map[string]string `json:"validations"`
}
