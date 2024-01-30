package user

import "inditilla/internal/repository/user"

type UserService interface {
}

type userService struct {
	userRepo user.UserRepo
}

func NewUserService(u user.UserRepo) *userService {
	return &userService{
		userRepo: u,
	}
}
