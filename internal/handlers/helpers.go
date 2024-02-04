package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"inditilla/internal/entity"
	"net/http"

	"github.com/go-playground/form/v4"
)

func (r *routes) decodePostForm(req *http.Request, dst any) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	if err := r.fd.Decode(dst, req.PostForm); err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (r *routes) logError(req *http.Request, err error) {
	r.l.Error("error: %v, request_method: %s, request_url: %s", err, req.Method, req.URL.String())
}

func (r *routes) sendResponse(w http.ResponseWriter, req *http.Request, status int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		r.sendErrorResponse(w, req, http.StatusInternalServerError, "Error marshaling response", "Response send")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}

func (r *routes) invalidAuthToken(w http.ResponseWriter, req *http.Request, location string) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	r.sendErrorResponse(w, req, http.StatusUnauthorized, "invalid or missing authentication token", location)
}

func (r *routes) unprocessableEntity(w http.ResponseWriter, req *http.Request, location string) {
	r.sendErrorResponse(w, req, http.StatusUnprocessableEntity, "invalid form fill", location)
}

func (r *routes) notFound(w http.ResponseWriter, req *http.Request, location string) {
	r.sendErrorResponse(w, req, http.StatusNotFound, "requested resource could not be found", location)
}

func (r *routes) badRequest(w http.ResponseWriter, req *http.Request, err error, location string) {
	r.sendErrorResponse(w, req, http.StatusBadRequest, err.Error(), location)
}

func (r *routes) serverError(w http.ResponseWriter, req *http.Request, err error, location string) {
	r.logError(req, err)

	r.sendErrorResponse(w, req, http.StatusInternalServerError, "server encountered an error and could not process your request", location)
}

func (r *routes) sendErrorResponse(w http.ResponseWriter, req *http.Request, status int, message string, location string) {
	errResp := entity.ErrorResponse{
		ResponseStatus: "fail",
		Code:           status,
		Message:        message,
		Location:       location,
	}

	jsonData, err := json.Marshal(errResp)
	if err != nil {
		r.logError(req, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.l.Error(fmt.Sprintf("message - %s, location - %s", message, location))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}
