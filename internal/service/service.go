package service

import (
	"inditilla/internal/repository"
	"inditilla/internal/service/user"
)

type Services struct {
	User user.UserService
}

func New(r *repository.Repositories) *Services {
	return &Services{
		User: user.NewUserService(r.User),
	}
}
