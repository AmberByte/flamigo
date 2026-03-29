package http

import (
	"encoding/json"
	"errors"
	stdhttp "net/http"

	flamigo "github.com/amberbyte/flamigo/core"
)

type ErrorEncoder func(w stdhttp.ResponseWriter, r *stdhttp.Request, err error)
type PayloadEncoder func(w stdhttp.ResponseWriter, r *stdhttp.Request, payload any)

func DefaultErrorEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
	status := stdhttp.StatusInternalServerError
	message := err.Error()

	var publicErr flamigo.PublicError
	if errors.As(err, &publicErr) {
		if publicErr.StatusCode() != 0 {
			status = publicErr.StatusCode()
		}
		message = publicErr.PublicError()
	}

	stdhttp.Error(w, message, status)
}

func DefaultPayloadEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, payload any) {
	if payload == nil {
		w.WriteHeader(stdhttp.StatusNoContent)
		return
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		DefaultErrorEncoder(w, r, flamigo.WrapError("encode http response: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stdhttp.StatusOK)
	_, _ = w.Write(raw)
}
