package entity

import (
	"inditilla/internal/service/validator"
	"time"
)

type UserEntity struct {
	FirstName      string
	LastName       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}

type UserSignupForm struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	validator.Validator
}

type UserLoginForm struct {
	Email    string
	Password string
	validator.Validator
}
