package user

import (
	"context"
	"errors"
	"inditilla/internal/entity"
	"inditilla/internal/repository/user"
	"net/http"
)

type UserService interface {
	SaveUser(context.Context, *entity.UserSignupForm) (int, int, error)
	// Authenticate(context.Context, *entity.UserLoginForm) (int, int, error)
}

type userService struct {
	userRepo user.UserRepo
}

func NewUserService(u user.UserRepo) *userService {
	return &userService{
		userRepo: u,
	}
}

func (us *userService) SaveUser(ctx context.Context, u *entity.UserSignupForm) (int, int, error) {
	if !isRightSignUp(u) {
		return 0, http.StatusUnprocessableEntity, nil
	}

	id, err := us.userRepo.SaveUser(ctx, *u)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicateEmail) {
			u.AddFieldError("email", "Email address is already in use")
			return 0, http.StatusUnprocessableEntity, err
		} else {
			return 0, http.StatusInternalServerError, err
		}
	}

	return id, http.StatusOK, nil
}
