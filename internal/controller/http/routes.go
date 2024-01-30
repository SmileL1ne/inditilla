package http

import (
	"inditilla/internal/service"
	"inditilla/pkg/logger"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type routes struct {
	l logger.ILogger
	s *service.Services
}

func NewRouter(logger logger.ILogger, services *service.Services) http.Handler {
	router := httprouter.New()

	r := &routes{
		l: logger,
		s: services,
	}

	router.HandlerFunc(http.MethodGet, "/v1/user/signup", r.userSignup)
	router.HandlerFunc(http.MethodPost, "/v1/user/signup", r.userSignupPost)
	router.HandlerFunc(http.MethodGet, "/v1/user/login", r.userLogin)
	router.HandlerFunc(http.MethodPost, "/v1/user/login", r.userLoginPost)

	return router
}
