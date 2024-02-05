package service

import (
	"inditilla/internal/data"
	"inditilla/internal/repository"
	"inditilla/internal/service/user"
)

type Services struct {
	User user.UserService
}

// New returns Services struct with all services initialized (only User service in this case)
func New(r *repository.Repositories, auth *user.Authorizer, tokenModel *data.TokenModel) *Services {
	return &Services{
		User: user.NewUserService(r.User, auth, tokenModel),
	}
}
