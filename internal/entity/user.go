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
	FirstName           string `form:"firstName"`
	LastName            string `form:"lastName"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type UserLoginForm struct {
	Email    string
	Password string
	validator.Validator
}
