package websocket_test

import (
	"errors"
	"testing"

	flamigo "github.com/amberbyte/flamigo/core"
	transportws "github.com/amberbyte/flamigo/transport/websocket"
	"github.com/stretchr/testify/assert"
)

type mockValidationError struct {
	err    error
	fields []flamigo.FieldError
}

func (e mockValidationError) Error() string {
	return e.err.Error()
}

func (e mockValidationError) Unwrap() error {
	return e.err
}

func (e mockValidationError) FieldErrors() []flamigo.FieldError {
	return e.fields
}

func TestEncodeSuccess(t *testing.T) {
	raw, err := transportws.EncodeSuccess("app::users:get", map[string]string{"id": "42"}, transportws.WithAckKey("abc"))
	assert.NoError(t, err)
	assert.JSONEq(t, `{"topic":"app::users:get","ack":"abc","payload":{"id":"42"}}`, string(raw))
}

func TestEncodeError(t *testing.T) {
	raw, err := transportws.EncodeError(flamigo.NewError("nope", flamigo.Public("Forbidden"), flamigo.Kind(flamigo.ErrorTypeForbidden)), transportws.WithAckKey("abc"))
	assert.NoError(t, err)
	assert.Contains(t, string(raw), `"topic":"error"`)
	assert.Contains(t, string(raw), `"ack":"abc"`)
	assert.Contains(t, string(raw), `"message":"Forbidden"`)
	assert.Contains(t, string(raw), `"type":"forbidden"`)
}

func TestEncodeWrappedPublicError(t *testing.T) {
	raw, err := transportws.EncodeError(flamigo.WrapError("wrapped: %w", flamigo.NewError("nope", flamigo.Public("Forbidden"), flamigo.Kind(flamigo.ErrorTypeForbidden))))
	assert.NoError(t, err)
	assert.Contains(t, string(raw), `"message":"Forbidden"`)
	assert.Contains(t, string(raw), `"type":"forbidden"`)
}

func TestEncodeErrorWithValidationFields(t *testing.T) {
	raw, err := transportws.EncodeError(mockValidationError{
		err: errors.New("invalid payload"),
		fields: []flamigo.FieldError{
			{Field: "token", Error: "token is required"},
		},
	})
	assert.NoError(t, err)
	assert.Contains(t, string(raw), `"fieldErrors":[{"field":"token","error":"token is required"}]`)
}
