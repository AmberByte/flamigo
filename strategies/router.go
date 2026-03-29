package strategies

import (
	"errors"
	"fmt"
	"strings"
)

type router[CTX Context] struct {
	registries map[string]Registry[CTX]
}

func validateNamespace(namespace string) error {
	if namespace == "" {
		return errors.New("namespace cannot be empty")
	}
	if strings.Contains(namespace, "::") {
		return fmt.Errorf("namespace %s must not contain namespace separator ::", namespace)
	}
	return nil
}

func splitAction(action string) (string, string, error) {
	namespace, localAction, ok := strings.Cut(action, "::")
	if !ok || namespace == "" || localAction == "" {
		return "", "", fmt.Errorf("action %s should be <namespace>::<action>", action)
	}
	return namespace, localAction, nil
}

func (r *router[CTX]) Register(namespace string, registry Registry[CTX]) error {
	if err := validateNamespace(namespace); err != nil {
		return fmt.Errorf("register router namespace: %w", err)
	}
	if _, exists := r.registries[namespace]; exists {
		return fmt.Errorf("register router namespace: registry already registered for %s", namespace)
	}
	r.registries[namespace] = registry
	return nil
}

func (r *router[CTX]) Invoke(ctx CTX) Result {
	namespace, action, err := splitAction(ctx.Request().Action())
	if err != nil {
		ctx.Response().SetError(err)
		return ctx.Response()
	}
	registry, ok := r.registries[namespace]
	if !ok {
		ctx.Response().SetError(fmt.Errorf("no registry found for namespace %s", namespace))
		return ctx.Response()
	}
	return registry.Invoke(action, ctx)
}

func NewRouter[CTX Context]() Router[CTX] {
	return &router[CTX]{
		registries: map[string]Registry[CTX]{},
	}
}
