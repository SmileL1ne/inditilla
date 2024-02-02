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
		// RETURN ERROR RESPONSE JSON
		r.l.Warn("client error - %v", err)
		return
	}

	_, status, err := r.s.User.SignUp(req.Context(), &form)
	if status != http.StatusOK {
		if status == http.StatusUnprocessableEntity {
			r.l.Error("invalid signup form fill")
			// RETURN ERROR RESPONSE JSON WITH 422 CODE AND WITH FORM DETAILS IN 'DETAILS' SECTION
		} else {
			r.l.Error("server error - %v", err)
			// RETURN ERROR RESPONSE JSON
		}

		return
	}

	// RETURN USER DTO
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully saved user"))
}

func (r *routes) userLogin(w http.ResponseWriter, req *http.Request) {
	var form entity.UserLoginForm

	err := r.decodePostForm(req, &form)
	if err != nil {
		// RETURN ERROR RESPONSE JSON
		r.l.Warn("client error - %v", err)
		return
	}

	token, status, err := r.s.User.SignIn(req.Context(), &form)
	if status != http.StatusOK {
		if status == http.StatusUnprocessableEntity {
			if errors.Is(err, entity.ErrInvalidCredentials) {
				form.AddNonFieldError("Email or password is incorrect")
			}
			r.l.Error("invalid login form fill")
			// RETURN ERROR RESPONSE JSON WITH 422 CODE AND WITH FORM DETAILS IN 'DETAILS' SECTION
		} else {
			r.l.Error("server error - %v", err)
			// RETURN ERROR RESPONSE JSON
		}

		return
	}

	// RETURN USER DTO WITH JWT TOKEN
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully saved user - " + token))
}

func (r *routes) userProfile(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm authorized!"))
}
