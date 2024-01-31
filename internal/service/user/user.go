package user

import (
	"inditilla/internal/entity"
	"inditilla/internal/repository/user"
)

type UserService interface {
	SaveUser(*entity.UserSignupForm) (int, int, error)
	Authenticate(*entity.UserLoginForm)
}

type userService struct {
	userRepo user.UserRepo
}

func NewUserService(u user.UserRepo) *userService {
	return &userService{
		userRepo: u,
	}
}
