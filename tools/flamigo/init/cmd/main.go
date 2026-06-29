package main

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strings"

	"init/internal/api"
	
	
	
)
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
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	slog.Info("starting init backend")
	injector := injection.NewInjector()

	for _, init := range initializers {
		err := validInitializer(init)
		if err != nil {
			slog.Error("verifying initializer", "initializer", getFunctionPackageName(init), "error", err)
			os.Exit(1)
		}
		err = injector.Invoke(init)
		if err != nil {
			slog.Error("initializing", "initializer", getFunctionPackageName(init), "error", err)
			os.Exit(1)
		}
	}
}
