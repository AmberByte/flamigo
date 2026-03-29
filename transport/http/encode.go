package http

import (
	"encoding/json"
	stdhttp "net/http"

	flamigo "github.com/amberbyte/flamigo/core"
)

type ErrorEncoder func(w stdhttp.ResponseWriter, r *stdhttp.Request, err error)
type PayloadEncoder func(w stdhttp.ResponseWriter, r *stdhttp.Request, payload any)

func DefaultErrorEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
	status := mapErrorTypeToStatus(flamigo.ResolveErrorType(err))
	message := flamigo.PublicMessage(err)

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

func mapErrorTypeToStatus(errorType flamigo.ErrorType) int {
	switch errorType {
	case flamigo.ErrorTypeBadRequest:
		return stdhttp.StatusBadRequest
	case flamigo.ErrorTypeUnauthorized:
		return stdhttp.StatusUnauthorized
	case flamigo.ErrorTypeForbidden:
		return stdhttp.StatusForbidden
	case flamigo.ErrorTypeNotFound:
		return stdhttp.StatusNotFound
	case flamigo.ErrorTypeConflict:
		return stdhttp.StatusConflict
	default:
		return stdhttp.StatusInternalServerError
	}
}
