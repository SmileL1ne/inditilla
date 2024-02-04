package entity

import (
	"inditilla/internal/service/validator"
	"time"
)

type UserEntity struct {
	Id             int
	FirstName      string
	LastName       string
	Email          string
	HashedPassword []byte
	CreatedAt      time.Time
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type SignupResponse struct {
	UserID int `json:"user_id"`
}

type UserProfileResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UserSignupForm struct {
	FirstName           string `form:"firstName"`
	LastName            string `form:"lastName"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
