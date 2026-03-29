package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	flamigo "github.com/amberbyte/flamigo/core"
	eventclient "github.com/amberbyte/flamigo/events/client"
	"github.com/go-playground/validator/v10"
)

type EncodeOption func(*responseBody)

type responseBody struct {
	Topic   string
	AckKey  string
	Payload any
}

func WithAckKey(ackKey string) EncodeOption {
	return func(r *responseBody) {
		r.AckKey = ackKey
	}
}

func FromCommand(cmd *Command) EncodeOption {
	return func(r *responseBody) {
		r.AckKey = cmd.AckKey()
	}
}

func marshalPayload(payload any) ([]byte, error) {
	switch v := payload.(type) {
	case eventclient.MessageMarshaller:
		return v.MarshalClientPayload()
	case json.Marshaler:
		return v.MarshalJSON()
	default:
		return json.Marshal(payload)
	}
}

func (r *responseBody) MarshalJSON() ([]byte, error) {
	rawPayload, err := marshalPayload(r.Payload)
	if err != nil {
		return nil, err
	}
	rawBody := struct {
		Topic   string          `json:"topic"`
		AckKey  string          `json:"ack,omitempty"`
		Payload json.RawMessage `json:"payload"`
	}{
		Topic:   r.Topic,
		AckKey:  r.AckKey,
		Payload: rawPayload,
	}
	return json.Marshal(rawBody)
}

func EncodeSuccess(topic string, payload any, opts ...EncodeOption) ([]byte, error) {
	body := &responseBody{
		Topic:   topic,
		Payload: payload,
	}
	for _, opt := range opts {
		opt(body)
	}
	return body.MarshalJSON()
}

type errorPayload struct {
	Message     string                `json:"message,omitempty"`
	Status      int                   `json:"status,omitempty"`
	FieldErrors []FieldErrorFormatted `json:"fieldErrors,omitempty"`
	Trace       any                   `json:"trace,omitempty"`
}

type FieldErrorFormatted struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func EncodeError(err error, opts ...EncodeOption) ([]byte, error) {
	publicErr := unwrapPublicError(err)
	payload := errorPayload{
		Message: publicErr.PublicError(),
		Status:  publicErr.StatusCode(),
		Trace:   err.Error(),
	}
	switch v := err.(type) {
	case validator.ValidationErrors:
		payload.FieldErrors = formatValidationError(v)
	case flamigo.PublicError:
		payload.Message = v.PublicError()
	}

	body := &responseBody{
		Topic:   "error",
		Payload: payload,
	}
	for _, opt := range opts {
		opt(body)
	}
	return body.MarshalJSON()
}

func unwrapPublicError(err error) flamigo.PublicError {
	unwrappedErr := errors.Unwrap(err)
	if err, ok := unwrappedErr.(flamigo.PublicError); ok {
		return err
	}
	if unwrappedErr == nil {
		return flamigo.WrapError("backend error: %w", err)
	}
	return unwrapPublicError(unwrappedErr)
}

func formatValidationError(err validator.ValidationErrors) []FieldErrorFormatted {
	formatted := make([]FieldErrorFormatted, 0, len(err))
	for _, e := range err {
		formatted = append(formatted, FieldErrorFormatted{
			Field: e.Field(),
			Error: formatFieldError(e),
		})
	}
	return formatted
}

func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of %s", fe.Field(), strings.ReplaceAll(fe.Param(), " ", ","))
	default:
		return fmt.Sprintf("%s failed for %s", fe.Field(), fe.Tag())
	}
}
