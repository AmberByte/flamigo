package strategies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterDispatchesByNamespace(t *testing.T) {
	appRegistry := NewRegistry[Context]()
	assert.NoError(t, appRegistry.Register("users:create", func(ctx Context) {
		ctx.Response().SetPayload("ok")
	}))

	router := NewRouter[Context]()
	assert.NoError(t, router.Register("app", appRegistry))

	ctx := NewContext(nil, "app::users:create", nil)
	result := router.Invoke(ctx)
	assert.NoError(t, result.Err())
	assert.Equal(t, "ok", result.Payload())
}

func TestRouterFailsOnUnknownNamespace(t *testing.T) {
	router := NewRouter[Context]()
	ctx := NewContext(nil, "admin::users:create", nil)
	result := router.Invoke(ctx)
	assert.ErrorContains(t, result.Err(), "no registry found for namespace admin")
}

func TestRouterFailsOnInvalidAction(t *testing.T) {
	router := NewRouter[Context]()
	ctx := NewContext(nil, "users:create", nil)
	result := router.Invoke(ctx)
	assert.ErrorContains(t, result.Err(), "should be <namespace>::<action>")
}

func TestRouterRejectsDuplicateNamespace(t *testing.T) {
	router := NewRouter[Context]()
	assert.NoError(t, router.Register("app", NewRegistry[Context]()))
	err := router.Register("app", NewRegistry[Context]())
	assert.ErrorContains(t, err, "registry already registered")
}
