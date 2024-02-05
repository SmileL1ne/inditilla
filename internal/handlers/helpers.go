package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"inditilla/internal/entity"
	"io"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (r *routes) readJSON(w http.ResponseWriter, req *http.Request, target interface{}) error {
	maxBytes := 1_048_576
	req.Body = http.MaxBytesReader(w, req.Body, int64(maxBytes))

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(target)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBodyLen *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON at character - %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field - %q", unmarshalTypeError.Field)
			} else {
				return fmt.Errorf("body contains incorrect JSON type at - %d", unmarshalTypeError.Offset)
			}
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unkown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unkown key - %s", fieldName)
		case errors.As(err, &maxBodyLen):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (r *routes) logError(req *http.Request, err error) {
	r.l.Error("error: %v, request_method: %s, request_url: %s", err, req.Method, req.URL.String())
}

func (r *routes) validateToken(token string) bool {
	return token != ""
}

func (r *routes) sendResponse(w http.ResponseWriter, req *http.Request, status int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		r.sendErrorResponse(w, req, http.StatusInternalServerError, "Error marshaling response", nil, "Response send")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(jsonData); err != nil {
		r.l.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (r *routes) invalidAuthToken(w http.ResponseWriter, req *http.Request, location string) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	r.sendErrorResponse(w, req, http.StatusUnauthorized, "invalid or missing authentication token", nil, location)
}

func (r *routes) editConflict(w http.ResponseWriter, req *http.Request, validations map[string]string, location string) {
	r.sendErrorResponse(w, req, http.StatusConflict, "unable to update the record due to an edit conflict, please try again", validations, location)
}

func (r *routes) unprocessableEntity(w http.ResponseWriter, req *http.Request, validations map[string]string, location string) {
	r.sendErrorResponse(w, req, http.StatusUnprocessableEntity, "invalid form fill", nil, location)
}

func (r *routes) notFound(w http.ResponseWriter, req *http.Request, location string) {
	r.sendErrorResponse(w, req, http.StatusNotFound, "requested resource could not be found", nil, location)
}

func (r *routes) badRequest(w http.ResponseWriter, req *http.Request, err error, location string) {
	r.sendErrorResponse(w, req, http.StatusBadRequest, err.Error(), nil, location)
}

func (r *routes) serverError(w http.ResponseWriter, req *http.Request, err error, location string) {
	r.logError(req, err)

	r.sendErrorResponse(w, req, http.StatusInternalServerError, "server encountered an error and could not process your request", nil, location)
}

func (r *routes) sendErrorResponse(w http.ResponseWriter, req *http.Request, status int, message string, validations map[string]string, location string) {
	errResp := entity.ErrorResponse{
		ResponseStatus: "fail",
		Code:           status,
		Message:        message,
		Location:       location,
	}

	if validations != nil {
		errResp.Validations = validations
	}

	jsonData, err := json.Marshal(errResp)
	if err != nil {
		r.logError(req, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// r.l.Error(fmt.Sprintf("message - %s, location - %s", message, location))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(jsonData); err != nil {
		r.l.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (r *routes) retrieveParamId(req *http.Request) string {
	req.URL.Path = httprouter.CleanPath(req.URL.Path)
	params := httprouter.ParamsFromContext(req.Context())
	return params.ByName("id")
}
