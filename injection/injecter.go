package injection

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Repository interface {
	AddInjectable(i any) error
}

type Provider interface {
	Execute(t any, args ...any) error
	ExecuteList(t []any, args ...any) error
}

type DependencyManager interface {
	Repository
	Provider
}

type ProviderRepository interface {
	Repository
	Provider
}

var (
	errInvalidExecutable = errors.New("provided executable is not a function")
)

func errInvalidInjectable(t reflect.Type) error {
	return fmt.Errorf("no injectable registered for type %s", t.String())
}

func errAlreadyRegistered(t reflect.Type) error {
	return fmt.Errorf("injectable of type %s is already registered", t.String())
}

var _ Repository = (*injecter)(nil)
var _ Provider = (*injecter)(nil)
var _ DependencyManager = (*injecter)(nil)

type injecter struct {
	injectables map[reflect.Type]reflect.Value
}

func (injector *injecter) AddInjectable(i any) error {
	t := reflect.TypeOf(i)
	if injector.injectables[t].IsValid() {
		return errAlreadyRegistered(t)
	}
	injector.injectables[t] = reflect.ValueOf(i)
	logrus.Debugf("Added injectable: %s", t.String())
	return nil
}

func (injector *injecter) getInjectable(t reflect.Type, args ...map[reflect.Type]reflect.Value) reflect.Value {
	result := findInjectable(injector.injectables, t)
	if result.IsValid() {
		return result
	}
	if len(args) > 0 {
		result = findInjectable(args[0], t)
	}
	return result
}

func (injector *injecter) Execute(t any, args ...any) error {
	tt := reflect.TypeOf(t)
	if tt.Kind() != reflect.Func {
		return errInvalidExecutable
	}
	parameters := []reflect.Value{}
	for i := 0; i < tt.NumIn(); i++ {
		injectable := injector.getInjectable(tt.In(i), transformArgs(args...))
		if !injectable.IsValid() {
			fmt.Printf("Injectable: %#v\n", tt.In(i))
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

func (injector *injecter) ExecuteList(t []any, args ...any) error {
	for _, f := range t {
		err := injector.Execute(f, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewDependencyInjecter() *injecter {
	injector := &injecter{
		injectables: map[reflect.Type]reflect.Value{},
	}
	// Add the injector itself so it can be injected
	injector.AddInjectable(injector)
	return injector
}
