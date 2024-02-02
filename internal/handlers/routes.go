package handlers

import (
	"inditilla/internal/service"
	"inditilla/pkg/logger"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
)

type routes struct {
	l  logger.ILogger
	s  *service.Services
	fd *form.Decoder
}

func NewRouter(logger logger.ILogger, services *service.Services, formDecoder *form.Decoder) http.Handler {
	router := httprouter.New()

	r := &routes{
		l:  logger,
		s:  services,
		fd: formDecoder,
	}

	router.HandlerFunc(http.MethodGet, "/v1/user/signup", r.userSignup)
	router.HandlerFunc(http.MethodPost, "/v1/user/signup", r.userSignupPost)
	router.HandlerFunc(http.MethodGet, "/v1/user/login", r.userLogin)
	router.HandlerFunc(http.MethodPost, "/v1/user/login", r.userLoginPost)

	return router
}