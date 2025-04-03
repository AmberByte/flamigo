package flamigo_test

import (
	"testing"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/stretchr/testify/assert"
)

func TestError_StatusCode(t *testing.T) {
	t.Run("Returns status code", func(t *testing.T) {
		err := flamigo.NewError("some error")
		assert.Equal(t, 500, err.StatusCode())
	})
	t.Run("Can overwrite status code", func(t *testing.T) {
		err := flamigo.NewError("some error", flamigo.StatusCode(400))
		assert.Equal(t, 400, err.StatusCode())
	})
}

func TestError_PublicError(t *testing.T) {
	t.Run("Returns public error", func(t *testing.T) {
		err := flamigo.NewError("some error", flamigo.Public("public error message"))
		assert.Equal(t, "public error message", err.PublicError())
	})
	t.Run("Returns inner error when no public message", func(t *testing.T) {
		err := flamigo.NewError("some error")
		assert.Equal(t, "some error", err.PublicError())
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
