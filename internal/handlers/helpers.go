package handlers

import (
	"encoding/json"
	"errors"
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

func (r *routes) sendResponse(w http.ResponseWriter, status int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		r.sendErrorResponse(w, http.StatusInternalServerError, "Error marshaling response", "Response send")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}

func (r *routes) sendErrorResponse(w http.ResponseWriter, status int, message string, location string) {
	errResp := entity.ErrorResponse{
		ResponseStatus: "fail",
		Code:           status,
		Message:        message,
		Location:       location,
	}

	jsonData, err := json.Marshal(errResp)
	if err != nil {
		http.Error(w, "Server could not process your request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}
