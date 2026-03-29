package strategies

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryAddingStrategyWithPrefix(t *testing.T) {
	registry := NewRegistry[Context]()
	assert.NoError(t, registry.Register("test", func(ctx Context) {}))
}

func TestAddingStrategyWithoutPrefix(t *testing.T) {
	registry := NewRegistry[Context]()
	err := registry.Register("test::nested", func(ctx Context) {})
	assert.ErrorContains(t, err, "must be local")
}

func TestAddingDuplicateStrategy(t *testing.T) {
	registry := NewRegistry[Context]()
	assert.NoError(t, registry.Register("test", func(ctx Context) {}))
	err := registry.Register("test", func(ctx Context) {})
	assert.ErrorContains(t, err, "strategy already registered")
}

func TestRegistryCallingStrategy(t *testing.T) {
	registry := NewRegistry[Context]()
	called := false
	assert.NoError(t, registry.Register("test", func(ctx Context) {
		called = true
	}))
	ctx := NewContext(nil, "test::test", map[string]interface{}{})
	result := registry.Invoke("test", ctx)
	assert.NoError(t, result.Err())
	assert.True(t, called)
}

func TestRegistryStrategyErrors(t *testing.T) {
	registry := NewRegistry[Context]()
	called := false
	assert.NoError(t, registry.Register("test", func(ctx Context) {
		called = true
		ctx.Response().SetError(errors.New("some error is here"))
	}))
	ctx := NewContext(nil, "test::test", map[string]interface{}{})
	result := registry.Invoke("test", ctx)
	assert.True(t, called, "should have been called")
	assert.ErrorContains(t, result.Err(), "some error is here", "should have returned error")
}

func TestRegistryCallingNonExistentStrategy(t *testing.T) {
	registry := NewRegistry[Context]()
	ctx := NewContext(nil, "test::test2", map[string]interface{}{})
	result := registry.Invoke("test2", ctx)
	assert.ErrorContains(t, result.Err(), "no handler found for test2")
}
