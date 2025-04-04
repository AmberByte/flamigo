package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"init/internal/api"
	
	
	
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
}

var initializers = []any{
	//------------  Core domains and packages
	
	flamigo.Init,
	//------------ Domains Infra
	
	// ----------- Domain Apps
	
	//------------ Initialuze APIs
	api.Init,
	
}

func wrapError(t reflect.Type, err error) error {
	return fmt.Errorf("initializer %s: %w", t.PkgPath(), err)

}

func validInitializer(init any) error {
	t := reflect.TypeOf(init)
	if t.Kind() != reflect.Func {
		return wrapError(t, fmt.Errorf("initializer must be a function, got %s", t.Kind()))
	}
	if t.NumOut() &lt; 1 {
		return wrapError(t, fmt.Errorf("initializer must return at least one value"))
	}

	return nil
}

func getFunctionPackageName(f interface{}) string {
	ptr := reflect.ValueOf(f).Pointer()
	fn := runtime.FuncForPC(ptr)
	if fn == nil {
		return "unknown"
	}

	fullName := fn.Name()
	return strings.Replace(fullName, "init", "", 1)
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	lgr.Info("Starting init backend")
	injector := injection.NewInjecter()

	for _, init := range initializers {
		err := validInitializer(init)
		if err != nil {
			lgr.Fatalf("verifying initializer (#%s): %s", getFunctionPackageName(init), err)
		}
		err = injector.Execute(init)
		if err != nil {

			lgr.Fatalf("initializing (%s): %s", getFunctionPackageName(init), err)
		}
	}
}
