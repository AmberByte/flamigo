package flamigo

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationError interface {
	error
	FieldErrors() []FieldError
}
