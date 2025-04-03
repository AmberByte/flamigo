package flamigo

import (
	"errors"
	"fmt"

	"github.com/amberbyte/flamigo/internal"
)

var _ error = (*Error)(nil)
var _ PublicError = (*Error)(nil)

type PublicError interface {
	PublicError() string
	StatusCode() int
}

type Error struct {
	innerError    error
	publicMessage string
	statusCode    int
}

func (e *Error) Error() string {
	return e.innerError.Error()
}

func (e *Error) PublicError() string {
	if e.publicMessage != "" {
		return e.publicMessage
	}

	if e.innerError != nil {
		return e.innerError.Error()
	}
	return "unknown error"
}

func (e *Error) Unwrap() error {
	return e.innerError
}

func (e *Error) StatusCode() int {
	return e.statusCode
}

type ErrorOpt = func(e *Error)

func StatusCode(code int) ErrorOpt {
	return func(e *Error) {
		e.statusCode = code
	}
}

func Public(message string) ErrorOpt {
	return func(e *Error) {
		e.publicMessage = message
	}
}

func WithPublicResponse(message string, code ...int) ErrorOpt {
	codeDefault := internal.ParseOptionalParam[int](code, 500)
	return func(e *Error) {
		e.publicMessage = message
		e.statusCode = codeDefault
	}
}

func NewError(message string, opts ...ErrorOpt) *Error {
	e := &Error{
		statusCode: 500,
		innerError: errors.New(message),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
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
