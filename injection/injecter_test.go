package injection

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type componentA struct {
	a string
}

type componentB struct {
	b string
}

type componentA2 struct {
	a string
}

func (c *componentB) GetB() string {
	return c.b
}

func (c *componentA) GetA() string {
	return c.a
}

func (c *componentA2) GetA() string {
	return c.a
}

func (c *componentA) GetB() string {
	return "bar"
}

type intA interface {
	GetA() string
}

type intB interface {
	GetB() string
}

func TestInjector_Register(t *testing.T) {
	t.Run("it adds an injectable (primitive)", func(t *testing.T) {
		injecter := NewInjector()
		err := injecter.Register("test")
		assert.NoError(t, err)
		actual, err := injecter.getInjectable(reflect.TypeOf("test"))
		assert.NoError(t, err)
		assert.Equal(t, reflect.ValueOf("test"), actual)
	})

	t.Run("it adds an injectable (pointer)", func(t *testing.T) {
		injecter := NewInjector()
		err := injecter.Register(&struct{}{})
		assert.NoError(t, err)
		actual, err := injecter.getInjectable(reflect.TypeOf(&struct{}{}))
		assert.NoError(t, err)
		assert.Equal(t, reflect.ValueOf(&struct{}{}), actual)
	})

	t.Run("it adds an injectable (multiple components)", func(t *testing.T) {
		injecter := NewInjector()
		err := injecter.Register(&componentA{})
		assert.NoError(t, err)

		b := &componentB{}
		err = injecter.Register(b)
		assert.NoError(t, err)
		actual, err := injecter.getInjectable(reflect.TypeOf(&componentB{}))
		assert.NoError(t, err)
		assert.Equal(t, reflect.ValueOf(b), actual)
	})
	t.Run("it can retrieve injectable (with interface)", func(t *testing.T) {
		injecter := NewInjector()
		a := &componentA{"hi"}
		err := injecter.Register(a)
		assert.NoError(t, err)

		actual, err := injecter.getInjectable(reflect.TypeOf((*intA)(nil)).Elem())
		assert.NoError(t, err)
		assert.Equal(t, reflect.ValueOf(a), actual)
	})
	t.Run("it rejects nil injectables", func(t *testing.T) {
		injecter := NewInjector()
		err := injecter.Register(nil)
		assert.ErrorIs(t, err, errNilInjectable)
	})
	t.Run("it rejects typed nil injectables", func(t *testing.T) {
		injecter := NewInjector()
		var dep *componentA
		err := injecter.Register(dep)
		assert.ErrorIs(t, err, errNilInjectable)
	})
}

func TestInjector_Invoke(t *testing.T) {
	t.Run("successfuly injects to a function", func(t *testing.T) {
		injecter := NewInjector()
		// Checking for properties to make sure the right one is returned
		a := &componentA{a: "Foo"}
		b := &componentB{b: "Bar"}
		injecter.Register(a)
		injecter.Register(b)

		err := injecter.Invoke(func(pa *componentA, pb *componentB) {
			assert.Equal(t, "Foo", pa.a)
			assert.Equal(t, "Bar", pb.b)
		})
		assert.NoError(t, err)
	})
	t.Run("successfuly injects with interface", func(t *testing.T) {
		injecter := NewInjector()
		// Checking for properties to make sure the right one is returned
		a := &componentA{a: "Foo"}
		b := &componentB{b: "Bar"}
		injecter.Register(a)
		injecter.Register(b)

		err := injecter.Invoke(func(pa intA, pb *componentB) {
			assert.Equal(t, "Foo", pa.GetA())
			assert.Equal(t, "Bar", pb.b)
		})
		assert.NoError(t, err)
	})
	t.Run("fails to inject to a non function", func(t *testing.T) {
		injecter := NewInjector()
		err := injecter.Invoke("hello")
		assert.Error(t, err)
	})
	t.Run("fails to inject to a function with missing injectable", func(t *testing.T) {
		injecter := NewInjector()
		injecter.Register(&componentA{})
		err := injecter.Invoke(func(pa *componentA, pb *componentB) {})
		assert.Error(t, err)
	})
	t.Run("fails to inject to a function with invalid injectable", func(t *testing.T) {
		injecter := NewInjector()
		injecter.Register(&componentA{})
		err := injecter.Invoke(func(pa string) {})
		assert.Error(t, err)
	})
	t.Run("fails to inject ambiguous interface implementations", func(t *testing.T) {
		injecter := NewInjector()
		injecter.Register(&componentA{})
		injecter.Register(&componentB{})
		err := injecter.Invoke(func(dep intB) {})
		assert.ErrorContains(t, err, "multiple injectables satisfy type")
	})
	t.Run("prefers explicit args over registered injectables", func(t *testing.T) {
		injecter := NewInjector()
		registered := &componentA{a: "registered"}
		override := &componentA{a: "override"}
		injecter.Register(registered)

		err := injecter.Invoke(func(dep *componentA) {
			assert.Equal(t, "override", dep.a)
		}, override)
		assert.NoError(t, err)
	})
	t.Run("prefers explicit interface args over registered injectables", func(t *testing.T) {
		injecter := NewInjector()
		registered := &componentA{a: "registered"}
		override := &componentA2{a: "override"}
		injecter.Register(registered)

		err := injecter.Invoke(func(dep intA) {
			assert.Equal(t, "override", dep.GetA())
		}, override)
		assert.NoError(t, err)
	})
	t.Run("fails when explicit args are ambiguous for an interface", func(t *testing.T) {
		injecter := NewInjector()
		err := injecter.Invoke(func(dep intA) {}, &componentA{}, &componentA2{})
		assert.ErrorContains(t, err, "multiple injectables satisfy type")
	})
	t.Run("fails to inject to a function with error return", func(t *testing.T) {
		injecter := NewInjector()
		injecter.Register(&componentA{})
		err := injecter.Invoke(func(pa *componentA) error {
			return assert.AnError
		})
		assert.Equal(t, assert.AnError, err)
	})
}
