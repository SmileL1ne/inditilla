package handlers

import (
	"errors"
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
