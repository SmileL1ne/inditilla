package handlers

import (
	"errors"
	"inditilla/internal/entity"
	"net/http"
)

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

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		user.Email = *input.Email
	}

	err = r.s.User.Update(req.Context(), &user)
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
}
