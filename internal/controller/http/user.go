package http

import "net/http"

func (r *routes) userSignup(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Nah, I'd win"))
}

func (r *routes) userSignupPost(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("user signup post"))
}

func (r *routes) userLogin(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("user login"))
}

func (r *routes) userLoginPost(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("user login post"))
}
