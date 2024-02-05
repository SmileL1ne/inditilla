package user

import (
	"context"
	"errors"
	"fmt"
	"inditilla/internal/data"
	"inditilla/internal/entity"
	"inditilla/internal/repository/user"
	"inditilla/internal/service/validator"
	"strconv"
	"time"
)

type UserService interface {
	SignUp(context.Context, *entity.UserSignupForm) (int, error)
	SignIn(context.Context, *entity.UserLoginForm) (string, error)
	Exists(context.Context, string) (bool, error)
	GetById(context.Context, string) (entity.UserEntity, error)
	Update(context.Context, *entity.UserEntity) error
}

type Authorizer struct {
	signingKey []byte
	deadline   time.Duration
}

func NewAuthorizer(signingKey []byte, deadline time.Duration) *Authorizer {
	return &Authorizer{
		signingKey: signingKey,
		deadline:   deadline,
	}
}

type userService struct {
	userRepo user.UserRepo
	auth     *Authorizer
	token    *data.TokenModel
}

func NewUserService(u user.UserRepo, auth *Authorizer, tokenModel *data.TokenModel) *userService {
	return &userService{
		userRepo: u,
		auth:     auth,
		token:    tokenModel,
	}
}

func (us *userService) SignUp(ctx context.Context, u *entity.UserSignupForm) (int, error) {
	if !isRightSignUp(u) {
		return 0, entity.ErrInvalidInputData
	}

	id, err := us.userRepo.SaveUser(ctx, *u)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicateEmail) {
			return 0, entity.ErrDuplicateEmail
		}
		return 0, err
	}

	return id, nil
}

func (us *userService) SignIn(ctx context.Context, u *entity.UserLoginForm) (string, error) {
	if !isRightLogin(u) {
		return "", entity.ErrInvalidInputData
	}

	user, err := us.userRepo.Authenticate(ctx, u.Email, u.Password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			return "", entity.ErrInvalidCredentials
		} else {
			return "", err
		}
	}

	token := us.token.New(user.Email, us.auth.deadline)

	tkn, err := token.SignedString(us.auth.signingKey)
	if err != nil {
		return "", fmt.Errorf("token signing error: %v", err)
	}

	return tkn, nil
}

func (us *userService) Exists(ctx context.Context, email string) (bool, error) {
	if !validator.Matches(email, EmailRX) {
		return false, nil
	}
	return us.userRepo.Exists(ctx, email)
}

func (us *userService) GetById(ctx context.Context, idStr string) (entity.UserEntity, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return entity.UserEntity{}, entity.ErrInvalidUserId
	}

	userEntity, err := us.userRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return entity.UserEntity{}, entity.ErrNoRecord
		}
		return entity.UserEntity{}, err
	}

	return userEntity, err
}

func (us *userService) Update(ctx context.Context, user *entity.UserEntity) error {
	if !isRightUser(user) {
		return entity.ErrInvalidInputData
	}

	return us.userRepo.Update(ctx, user)
}
