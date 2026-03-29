package injection

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Registrar interface {
	// Register adds an injectable to the container.
	Register(i any) error
}

type Invoker interface {
	// Invoke executes a function. Parameters are injected if candidates are known
	//
	// If the function returns error or injectable is not found, it will be returned.
	Invoke(t any, args ...any) error

	// InvokeAll executes a list of functions. Parameters are injected if candidates are known
	//
	// If any function returns error or injectable is not found, it will be returned.
	// When a function fails further functions will not be executed.
	InvokeAll(t []any, args ...any) error
}

type Container interface {
	Registrar
	Invoker
}

var (
	errInvalidExecutable = errors.New("provided executable is not a function")
)

func errInvalidInjectable(t reflect.Type) error {
	return fmt.Errorf("no injectable registered for type %s", t.String())
}

func errAmbiguousInjectable(t reflect.Type, matches string) error {
	return fmt.Errorf("multiple injectables satisfy type %s: %s", t.String(), matches)
}

func errAlreadyRegistered(t reflect.Type) error {
	return fmt.Errorf("injectable of type %s is already registered", t.String())
}

var _ Registrar = (*Injector)(nil)
var _ Invoker = (*Injector)(nil)
var _ Container = (*Injector)(nil)

type Injector struct {
	injectables map[reflect.Type]reflect.Value
}

func (injector *Injector) Register(i any) error {
	t := reflect.TypeOf(i)
	if injector.injectables[t].IsValid() {
		return errAlreadyRegistered(t)
	}
	injector.injectables[t] = reflect.ValueOf(i)
	logrus.Debugf("Added injectable: %s", t.String())
	return nil
}

func (injector *Injector) getInjectable(t reflect.Type, args ...map[reflect.Type]reflect.Value) (reflect.Value, error) {
	if len(args) > 0 {
		result, err := findInjectable(args[0], t)
		if err != nil || result.IsValid() {
			return result, err
		}
	}
	return findInjectable(injector.injectables, t)
}

func (injector *Injector) Invoke(t any, args ...any) error {
	tt := reflect.TypeOf(t)
	if tt.Kind() != reflect.Func {
		return errInvalidExecutable
	}
	argValues := transformArgs(args...)
	parameters := []reflect.Value{}
	for i := 0; i < tt.NumIn(); i++ {
		injectable, err := injector.getInjectable(tt.In(i), argValues)
		if err != nil {
			return err
		}
		if !injectable.IsValid() {
			return errInvalidInjectable(tt.In(i))
		}
		parameters = append(parameters, injectable)
	}
	v := reflect.ValueOf(t)
	results := v.Call(parameters)
	for _, r := range results {
		if err, ok := r.Interface().(error); ok {
			return err
		}
	}
	return nil
}

func (injector *Injector) InvokeAll(t []any, args ...any) error {
	for _, f := range t {
		err := injector.Invoke(f, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewInjector() *Injector {
	injector := &Injector{
		injectables: map[reflect.Type]reflect.Value{},
	}
	// Add the injector itself so it can be injected
	injector.Register(injector)
	return injector
}
