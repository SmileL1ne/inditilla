package handlers

import (
	"inditilla/internal/entity"
	"net/http"
)

func (r *routes) userSignup(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Nah, I'd win"))
}

func (r *routes) userSignupPost(w http.ResponseWriter, req *http.Request) {
	var form entity.UserSignupForm

	err := r.decodePostForm(req, &form)
	if err != nil {
		// RETURN ERROR RESPONSE JSON
		r.l.Warn("client error - %v", err)
		return
	}

	_, status, err := r.s.User.SaveUser(req.Context(), &form)
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

	// RETURN USER ENTITY
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully saved user"))
}

func (r *routes) userLogin(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("user login"))
}

func (r *routes) userLoginPost(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("user login post"))
}
