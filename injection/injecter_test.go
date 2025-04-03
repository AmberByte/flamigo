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

func (c *componentB) GetB() string {
	return c.b
}

func (c *componentA) GetA() string {
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

func TestInjecter_AddInjectable(t *testing.T) {
	t.Run("it adds an injectable (primitive)", func(t *testing.T) {
		injecter := NewDependencyInjecter()
		err := injecter.AddInjectable("test")
		assert.NoError(t, err)
		assert.Equal(t, reflect.ValueOf("test"), injecter.getInjectable(reflect.TypeOf("test")))
	})

	t.Run("it adds an injectable (pointer)", func(t *testing.T) {
		injecter := NewDependencyInjecter()
		err := injecter.AddInjectable(&struct{}{})
		assert.NoError(t, err)
		assert.Equal(t, reflect.ValueOf(&struct{}{}), injecter.getInjectable(reflect.TypeOf(&struct{}{})))
	})

	t.Run("it adds an injectable (multiple components)", func(t *testing.T) {
		injecter := NewDependencyInjecter()
		err := injecter.AddInjectable(&componentA{})
		assert.NoError(t, err)

		b := &componentB{}
		err = injecter.AddInjectable(b)
		assert.NoError(t, err)
		assert.Equal(t, reflect.ValueOf(b), injecter.getInjectable(reflect.TypeOf(&componentB{})))
	})
	t.Run("it can retrieve injectable (with interface)", func(t *testing.T) {
		injecter := NewDependencyInjecter()
		a := &componentA{"hi"}
		err := injecter.AddInjectable(a)
		assert.NoError(t, err)

		assert.Equal(t, reflect.ValueOf(a), injecter.getInjectable(reflect.TypeOf(intA(&componentA{}))))
	})
}

func TestInjecter_InjectToFunction(t *testing.T) {
	t.Run("successfuly injects to a function", func(t *testing.T) {
		injecter := NewDependencyInjecter()
		// Checking for properties to make sure the right one is returned
		a := &componentA{a: "Foo"}
		b := &componentB{b: "Bar"}
		injecter.AddInjectable(a)
		injecter.AddInjectable(b)

		err := injecter.Execute(func(pa *componentA, pb *componentB) {
			assert.Equal(t, "Foo", pa.a)
			assert.Equal(t, "Bar", pb.b)
		})
		assert.NoError(t, err)
	})
	t.Run("successfuly injects with interface", func(t *testing.T) {
		injecter := NewDependencyInjecter()
		// Checking for properties to make sure the right one is returned
		a := &componentA{a: "Foo"}
		b := &componentB{b: "Bar"}
		injecter.AddInjectable(a)
		injecter.AddInjectable(b)

		err := injecter.Execute(func(pa intA, pb *componentB) {
			assert.Equal(t, "Foo", pa.GetA())
			assert.Equal(t, "Bar", pb.b)
		})
		assert.NoError(t, err)
	})
	t.Run("fails to inject to a non function", func(t *testing.T) {
		injecter := NewDependencyInjecter()
		err := injecter.Execute("hello")
		assert.Error(t, err)
	})
	t.Run("fails to inject to a function with missing injectable", func(t *testing.T) {
		injecter := NewDependencyInjecter()
		injecter.AddInjectable(&componentA{})
		err := injecter.Execute(func(pa *componentA, pb *componentB) {})
		assert.Error(t, err)
	})
	t.Run("fails to inject to a function with invalid injectable", func(t *testing.T) {
		injecter := NewDependencyInjecter()
		injecter.AddInjectable(&componentA{})
		err := injecter.Execute(func(pa string) {})
		assert.Error(t, err)
	})
	t.Run("fails to inject to a function with error return", func(t *testing.T) {
		injecter := NewDependencyInjecter()
		injecter.AddInjectable(&componentA{})
		err := injecter.Execute(func(pa *componentA) error {
			return assert.AnError
		})
		assert.Equal(t, assert.AnError, err)
	})
}
