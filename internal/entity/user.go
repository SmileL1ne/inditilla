package entity

import (
	"inditilla/internal/service/validator"
	"time"
)

type UserEntity struct {
	Id                  int       `json:"id"`
	FirstName           string    `json:"firstName"`
	LastName            string    `json:"lastName"`
	Email               string    `json:"email"`
	Password            string    `json:"-"`
	CreatedAt           time.Time `json:"createdAt"`
	validator.Validator `json:"-"`
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
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	Email               string `json:"email"`
	Password            string `json:"password"`
	validator.Validator `json:"-"`
}

type UserLoginForm struct {
	Email               string `json:"email"`
	Password            string `json:"password"`
	validator.Validator `json:"-"`
}
