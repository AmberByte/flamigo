package strategies

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryAddingStrategyWithPrefix(t *testing.T) {
	registry := NewRegistry[Context]("test")
	assert.NotPanics(t, func() { registry.Register("test::test", func(ctx Context) {}) })
}

func TestAddingStrategyWithoutPrefix(t *testing.T) {
	registry := NewRegistry[Context]("test")
	assert.Panics(t, func() { registry.Register("test", func(ctx Context) {}) })
}

func TestRegistryCallingStrategy(t *testing.T) {
	registry := NewRegistry[Context]("test")
	called := false
	registry.Register("test::test", func(ctx Context) {
		called = true
	})
	ctx := NewContext(nil, "test::test", map[string]interface{}{})
	result := registry.Use(ctx)
	assert.NoError(t, result.Err())
	assert.True(t, called)
}

func TestRegistryStrategyErrors(t *testing.T) {
	registry := NewRegistry[Context]("test")
	called := false
	registry.Register("test::test", func(ctx Context) {
		called = true
		ctx.Response().SetError(errors.New("some error is here"))
	})
	ctx := NewContext(nil, "test::test", map[string]interface{}{})
	result := registry.Use(ctx)
	assert.True(t, called, "should have been called")
	assert.ErrorContains(t, result.Err(), "some error is here", "should have returned error")
}

func TestRegistryCallingNonExistentStrategy(t *testing.T) {
	registry := NewRegistry[Context]("test")
	ctx := NewContext(nil, "test::test2", map[string]interface{}{})
	result := registry.Use(ctx)
	assert.ErrorContains(t, result.Err(), "no handler found for test::test2")
}

func TestRegistryCallingWithWrongPrefix(t *testing.T) {
	registry := NewRegistry[Context]("test")
	ctx := NewContext(nil, "test2:test", map[string]interface{}{})
	result := registry.Use(ctx)
	assert.ErrorContains(t, result.Err(), "strategy test2:test should be test::test2:test")
}
