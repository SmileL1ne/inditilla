package handlers

import (
	"errors"
	"inditilla/internal/entity"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (r *routes) userSignup(w http.ResponseWriter, req *http.Request) {
	var form entity.UserSignupForm

	err := r.decodePostForm(req, &form)
	if err != nil {
		r.badRequest(w, req, err, "User signup")
		return
	}

	id, err := r.s.User.SignUp(req.Context(), &form)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormFill):
			r.unprocessableEntity(w, req, "User signup")
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
	var form entity.UserLoginForm

	err := r.decodePostForm(req, &form)
	if err != nil {
		r.badRequest(w, req, err, "User login")
		return
	}

	token, err := r.s.User.SignIn(req.Context(), &form)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormFill):
			r.unprocessableEntity(w, req, "User login")
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
	req.URL.Path = httprouter.CleanPath(req.URL.Path)
	params := httprouter.ParamsFromContext(req.Context())
	id := params.ByName("id")

	user, err := r.s.User.GetById(req.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoRecord):
			r.notFound(w, req, "User profile")
		case errors.Is(err, entity.ErrInvalidUserId):
			r.badRequest(w, req, err, "User profile")
		default:
			r.serverError(w, req, err, "User profile")
		}

		return
	}

	r.sendResponse(w, req, http.StatusOK, user)
}
