package websocket

import (
	"encoding/json"
	"errors"

	flamigo "github.com/amberbyte/flamigo/core"
	eventclient "github.com/amberbyte/flamigo/events/client"
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
	Message     string              `json:"message,omitempty"`
	Status      int                 `json:"status,omitempty"`
	FieldErrors []flamigo.FieldError `json:"fieldErrors,omitempty"`
	Trace       any                 `json:"trace,omitempty"`
}

func EncodeError(err error, opts ...EncodeOption) ([]byte, error) {
	publicErr := unwrapPublicError(err)
	payload := errorPayload{
		Message: publicErr.PublicError(),
		Status:  publicErr.StatusCode(),
		Trace:   err.Error(),
	}
	var validationErr flamigo.ValidationError
	if errors.As(err, &validationErr) {
		payload.FieldErrors = validationErr.FieldErrors()
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
	var fallback flamigo.PublicError
	for currentErr := err; currentErr != nil; currentErr = errors.Unwrap(currentErr) {
		publicErr, ok := currentErr.(flamigo.PublicError)
		if !ok {
			continue
		}
		if publicErr.StatusCode() != 0 {
			return publicErr
		}
		if fallback == nil {
			fallback = publicErr
		}
	}
	if fallback != nil {
		return fallback
	}
	return flamigo.WrapError("backend error: %w", err)
}
