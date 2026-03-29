package flamigo

import (
	"errors"
	"fmt"
)

var _ error = (*Error)(nil)
var _ PublicError = (*Error)(nil)

type ErrorType string

const (
	ErrorTypeBadRequest   ErrorType = "bad_request"
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	ErrorTypeForbidden    ErrorType = "forbidden"
	ErrorTypeNotFound     ErrorType = "not_found"
	ErrorTypeConflict     ErrorType = "conflict"
	ErrorTypeServerError  ErrorType = "server_error"
)

type PublicError interface {
	// PublicMessage returns the public error message
	// If no public message is set, it returns the inner error message
	// If no inner error is set, it returns "unknown error"
	PublicMessage() string
	Type() ErrorType
}

type Error struct {
	innerError    error
	publicMessage string
	errorType     ErrorType
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.innerError.Error()
}

// PublicMessage returns the public error message
// If no public message is set, it returns the inner error message
// If no inner error is set, it returns "unknown error"
func (e *Error) PublicMessage() string {
	if e.publicMessage != "" {
		return e.publicMessage
	}

	if e.innerError != nil {
		return e.innerError.Error()
	}
	return "unknown error"
}

// Unwraps the internal error
//
// If there is no inner error, it returns nil
func (e *Error) Unwrap() error {
	return e.innerError
}

func (e *Error) Type() ErrorType {
	if e.errorType != "" {
		return e.errorType
	}
	return ErrorTypeServerError
}

type ErrorOpt = func(e *Error)

// Kind sets the transport-agnostic error type.
func Kind(errorType ErrorType) ErrorOpt {
	return func(e *Error) {
		e.errorType = errorType
	}
}

// Public sets the public message
func Public(message string) ErrorOpt {
	return func(e *Error) {
		e.publicMessage = message
	}
}

// WithPublicResponse sets the public message and optional generic error type.
func WithPublicResponse(message string, errorType ...ErrorType) ErrorOpt {
	return func(e *Error) {
		e.publicMessage = message
		if len(errorType) > 0 {
			e.errorType = errorType[0]
		}
	}
}

func NewError(message string, opts ...ErrorOpt) *Error {
	e := &Error{
		innerError: errors.New(message),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func PublicMessage(err error) string {
	if err == nil {
		return "unknown error"
	}

	var fallback string
	for currentErr := err; currentErr != nil; currentErr = errors.Unwrap(currentErr) {
		switch e := currentErr.(type) {
		case *Error:
			if e.publicMessage != "" {
				return e.publicMessage
			}
			if fallback == "" {
				fallback = e.PublicMessage()
			}
		case PublicError:
			if fallback == "" {
				fallback = e.PublicMessage()
			}
		}
	}

	if fallback != "" {
		return fallback
	}

	return err.Error()
}

func ResolveErrorType(err error) ErrorType {
	if err == nil {
		return ErrorTypeServerError
	}

	var fallback ErrorType
	for currentErr := err; currentErr != nil; currentErr = errors.Unwrap(currentErr) {
		switch e := currentErr.(type) {
		case *Error:
			if e.errorType != "" {
				return e.errorType
			}
			if fallback == "" {
				fallback = e.Type()
			}
		case PublicError:
			if fallback == "" {
				fallback = e.Type()
			}
		}
	}

	if fallback != "" {
		return fallback
	}

	return ErrorTypeServerError
}

func WrapError(message string, err error, opts ...ErrorOpt) *Error {
	e := &Error{
		innerError: fmt.Errorf(message, err),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}
