package handlers

import (
	"errors"
	"fmt"
	"inditilla/internal/entity"
	"net/http"
	"time"
)

var timeFormat = "2006-01-02 15:04:05"

func (r *routes) userSignup(w http.ResponseWriter, req *http.Request) {
	var userSignupForm entity.UserSignupForm

	err := r.readJSON(w, req, &userSignupForm)
	if err != nil {
		r.badRequest(w, req, err, "User signup")
		return
	}

	id, err := r.s.User.SignUp(req.Context(), &userSignupForm)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidInputData):
			r.unprocessableEntity(w, req, userSignupForm.Validator.FieldErrors, "User signup")
		case errors.Is(err, entity.ErrDuplicateEmail):
			r.badRequest(w, req, err, "User signup")
		default:
			r.serverError(w, req, err, "User singup")
		}

		return
	}

	signupResp := entity.SignupResponse{
		UserID: id,
	}

	r.sendResponse(w, req, http.StatusOK, signupResp)

	// Log new user sign up
	r.l.Info("new user with id '%d' signed up at %s", id, time.Now().Format(timeFormat))
}

func (r *routes) userLogin(w http.ResponseWriter, req *http.Request) {
	var userLoginForm entity.UserLoginForm

	err := r.readJSON(w, req, &userLoginForm)
	if err != nil {
		r.badRequest(w, req, err, "User login")
		return
	}

	token, err := r.s.User.SignIn(req.Context(), &userLoginForm)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidInputData):
			r.unprocessableEntity(w, req, userLoginForm.Validator.FieldErrors, "User login")
		case errors.Is(err, entity.ErrInvalidCredentials):
			r.badRequest(w, req, err, "User login")
		default:
			r.serverError(w, req, err, "User login")
		}

		return
	}

	loginResp := entity.LoginResponse{
		AccessToken: token,
	}

	r.sendResponse(w, req, http.StatusCreated, loginResp)

	// Log user log in
	r.l.Info("user with email '%s' logged at %s", userLoginForm.Email, time.Now().Format(timeFormat))
}

func (r *routes) userProfile(w http.ResponseWriter, req *http.Request) {
	id := r.retrieveParamId(req)

	user, err := r.s.User.GetById(req.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoRecord):
			r.notFound(w, req, "User profile")
		case errors.Is(err, entity.ErrInvalidUserId):
			r.notFound(w, req, "User profile")
		default:
			r.serverError(w, req, err, "User profile")
		}

		return
	}

	userProfile := entity.UserProfileResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	r.sendResponse(w, req, http.StatusOK, userProfile)
}

func (r *routes) userUpdate(w http.ResponseWriter, req *http.Request) {
	id := r.retrieveParamId(req)

	user, err := r.s.User.GetById(req.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoRecord):
			r.notFound(w, req, "User update")
		case errors.Is(err, entity.ErrInvalidUserId):
			r.notFound(w, req, "User update")
		default:
			r.serverError(w, req, err, "User update")
		}

		return
	}

	var input struct {
		FirstName *string `json:"firstName"`
		LastName  *string `json:"lastName"`
		Email     *string `json:"email"`
		Password  *string `json:"password"`
	}

	err = r.readJSON(w, req, &input)
	if err != nil {
		r.badRequest(w, req, err, "User update")
		return
	}

	updatedFieldsLog := []string{}
	isPasswordChanged := false

	if input.FirstName != nil {
		updatedFieldsLog = append(updatedFieldsLog, fmt.Sprintf("firstName:'%s'->'%s'", user.FirstName, *input.FirstName))
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		updatedFieldsLog = append(updatedFieldsLog, fmt.Sprintf("lastName:'%s'->'%s'", user.LastName, *input.LastName))
		user.LastName = *input.LastName
	}
	if input.Email != nil {
		updatedFieldsLog = append(updatedFieldsLog, fmt.Sprintf("email:'%s'->'%s'", user.Email, *input.Email))
		user.Email = *input.Email
	}
	if input.Password != nil {
		isPasswordChanged = true
		updatedFieldsLog = append(updatedFieldsLog, "updated password")
		user.Password = *input.Password
	}

	err = r.s.User.Update(req.Context(), &user, isPasswordChanged)
	if err != nil {
		if errors.Is(err, entity.ErrEditConflict) {
			r.editConflict(w, req, user.FieldErrors, "User update")
			return
		}
		r.serverError(w, req, err, "User update")
	}

	userProfile := entity.UserProfileResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	r.sendResponse(w, req, http.StatusOK, userProfile)

	// Log user profile changes
	r.l.Info("user with id '%d' made next changes at %s in profile page: %v", user.Id, time.Now().Format(timeFormat), updatedFieldsLog)
}
