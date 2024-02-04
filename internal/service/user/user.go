package user

import (
	"context"
	"errors"
	"fmt"
	"inditilla/internal/entity"
	"inditilla/internal/repository/user"
	"inditilla/internal/service/validator"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type UserService interface {
	SignUp(context.Context, *entity.UserSignupForm) (int, int, error)
	SignIn(context.Context, *entity.UserLoginForm) (string, int, error)
	Exists(context.Context, string) (bool, error)
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
}

func NewUserService(u user.UserRepo, auth *Authorizer) *userService {
	return &userService{
		userRepo: u,
		auth:     auth,
	}
}

func (us *userService) SignUp(ctx context.Context, u *entity.UserSignupForm) (int, int, error) {
	if !isRightSignUp(u) {
		return 0, http.StatusUnprocessableEntity, nil
	}

	id, err := us.userRepo.SaveUser(ctx, *u)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicateEmail) {
			return 0, http.StatusBadRequest, err
		}
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusOK, nil
}

func (us *userService) SignIn(ctx context.Context, u *entity.UserLoginForm) (string, int, error) {
	if !isRightLogin(u) {
		return "", http.StatusUnprocessableEntity, nil
	}

	user, err := us.userRepo.Authenticate(ctx, u.Email, u.Password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			return "", http.StatusUnprocessableEntity, nil
		} else {
			return "", http.StatusInternalServerError, err
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(us.auth.deadline)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Email: user.Email,
	})

	tkn, err := token.SignedString(us.auth.signingKey)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("token signing error: %v", err)
	}
	return tkn, http.StatusOK, nil
}

func (us *userService) Exists(ctx context.Context, email string) (bool, error) {
	if !validator.Matches(email, EmailRX) {
		return false, nil
	}
	return us.userRepo.Exists(ctx, email)
}
