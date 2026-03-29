package flamigo_test

import (
	"testing"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/stretchr/testify/assert"
)

func TestError_ErrorType(t *testing.T) {
	t.Run("Returns default error type", func(t *testing.T) {
		err := flamigo.NewError("some error")
		assert.Equal(t, flamigo.ErrorTypeServerError, err.Type())
	})
	t.Run("Can overwrite error type", func(t *testing.T) {
		err := flamigo.NewError("some error", flamigo.Kind(flamigo.ErrorTypeBadRequest))
		assert.Equal(t, flamigo.ErrorTypeBadRequest, err.Type())
	})
}

func TestError_PublicError(t *testing.T) {
	t.Run("Returns public error", func(t *testing.T) {
		err := flamigo.NewError("some error", flamigo.Public("public error message"))
		assert.Equal(t, "public error message", err.PublicMessage())
	})
	t.Run("Returns inner error when no public message", func(t *testing.T) {
		err := flamigo.NewError("some error")
		assert.Equal(t, "some error", err.PublicMessage())
	})
}

func TestPublicMessage(t *testing.T) {
	t.Run("prefers explicit wrapped public message", func(t *testing.T) {
		err := flamigo.WrapError("listing messages: %w",
			flamigo.NewError("not found", flamigo.Public("Messages not found"), flamigo.Kind(flamigo.ErrorTypeNotFound)),
		)
		assert.Equal(t, "Messages not found", flamigo.PublicMessage(err))
	})
}

func TestResolveErrorType(t *testing.T) {
	t.Run("prefers explicit wrapped error type", func(t *testing.T) {
		err := flamigo.WrapError("listing messages: %w",
			flamigo.NewError("not found", flamigo.Public("Messages not found"), flamigo.Kind(flamigo.ErrorTypeNotFound)),
		)
		assert.Equal(t, flamigo.ErrorTypeNotFound, flamigo.ResolveErrorType(err))
	})
}

func TestNewError(t *testing.T) {
	t.Run("Sets error message", func(t *testing.T) {
		err := flamigo.NewError("some error")
		assert.Equal(t, "some error", err.Error())
	})
}

func TestWrapError(t *testing.T) {
	t.Run("Wraps error", func(t *testing.T) {
		err := flamigo.WrapError("this is my message %w", assert.AnError)
		assert.ErrorContains(t, err, "this is my message")
		assert.ErrorContains(t, err, assert.AnError.Error())
	})
}
