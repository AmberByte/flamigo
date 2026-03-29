package strategies

import (
	"errors"
	"fmt"
	"strings"
)

type registry[CTX Context] struct {
	strategies map[string]Strategy[CTX]
}

func validateAction(action string) error {
	if action == "" {
		return errors.New("action cannot be empty")
	}
	if strings.Contains(action, "::") {
		return fmt.Errorf("action %s must be local and must not contain namespace separator ::", action)
	}
	return nil
}

func (r *registry[CTX]) Register(action string, fn Strategy[CTX]) error {
	if err := validateAction(action); err != nil {
		return fmt.Errorf("register strategy: %w", err)
	}
	if _, exists := r.strategies[action]; exists {
		return fmt.Errorf("register strategy: strategy already registered for %s", action)
	}
	r.strategies[action] = fn
	return nil
}

func (r *registry[CTX]) Invoke(action string, ctx CTX) Result {
	fn, ok := r.strategies[action]
	if !ok {
		ctx.Response().SetError(errors.New("no handler found for " + action))
		return ctx.Response()
	}
	fn(ctx)
	return ctx.Response()
}

func NewRegistry[CTX Context]() Registry[CTX] {
	return &registry[CTX]{
		strategies: map[string]Strategy[CTX]{},
	}
}
