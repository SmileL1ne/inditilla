package handlers

import (
	"errors"
	"inditilla/internal/entity"
	"net/http"
)

func (r *routes) userSignup(w http.ResponseWriter, req *http.Request) {
	var form entity.UserSignupForm

	err := r.decodePostForm(req, &form)
	if err != nil {
		r.l.Warn("client error - %v", err)
		r.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", "User signup")
		return
	}

	id, status, err := r.s.User.SignUp(req.Context(), &form)
	if status != http.StatusOK {
		if status == http.StatusUnprocessableEntity {
			r.l.Error("invalid user signup form fill")
			r.sendErrorResponse(w, status, "Invalid user signup form fill", "User signup")
		} else if errors.Is(err, entity.ErrDuplicateEmail) {
			r.l.Error("user with email '%s' already exists", form.Email)
			r.sendErrorResponse(w, http.StatusBadRequest, "Email already in use", "User signup")
		} else {
			r.l.Error("server error - %v", err)
			r.sendErrorResponse(w, http.StatusInternalServerError, err.Error(), "User signup")
		}

		return
	}

	signupResp := entity.SignupResponse{
		UserID: id,
	}

	r.sendResponse(w, http.StatusOK, signupResp)
}

func (r *routes) userLogin(w http.ResponseWriter, req *http.Request) {
	var form entity.UserLoginForm

	err := r.decodePostForm(req, &form)
	if err != nil {
		r.l.Warn("client error - %v", err)
		r.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", "User login")
		return
	}

	token, status, err := r.s.User.SignIn(req.Context(), &form)
	if status != http.StatusOK {
		if status == http.StatusUnprocessableEntity {
			if errors.Is(err, entity.ErrInvalidCredentials) {
				r.l.Error("email or password is incorrect")
				r.sendErrorResponse(w, http.StatusBadRequest, "Email or password is incorrect", "User login")
				form.AddNonFieldError("Email or password is incorrect")
			}
			r.l.Error("invalid login form fill")
			r.sendErrorResponse(w, http.StatusUnprocessableEntity, "Invalid login form fill", "User login")
		} else {
			r.l.Error("server error - %v", err)
			r.sendErrorResponse(w, http.StatusInternalServerError, err.Error(), "User login")
		}

		return
	}

	loginResp := entity.LoginResponse{
		AccessToken: token,
	}

	r.sendResponse(w, http.StatusOK, loginResp)
}

func (r *routes) userProfile(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm authorized!"))
}
