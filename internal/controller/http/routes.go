package http

import (
	"inditilla/pkg/logger"
	"net/http"
)

type router struct {
	logger logger.ILogger
	// add service here
}

func NewRouter(logger logger.ILogger, s *service.Service) http.Handler {

}
