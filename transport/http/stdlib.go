package http

import (
	"fmt"
	stdhttp "net/http"
	"strings"
)

var _ Registrar = (*ServeMuxRegistrar)(nil)

type ServeMuxRegistrar struct {
	mux        *stdhttp.ServeMux
	dispatcher *Dispatcher
}

func NewServeMuxRegistrar(mux *stdhttp.ServeMux, dispatcher *Dispatcher) *ServeMuxRegistrar {
	return &ServeMuxRegistrar{
		mux:        mux,
		dispatcher: dispatcher,
	}
}

func (r *ServeMuxRegistrar) Handle(method string, path string, action string) (err error) {
	if r.mux == nil {
		return fmt.Errorf("register http route %s %s: serve mux is nil", method, path)
	}
	if r.dispatcher == nil {
		return fmt.Errorf("register http route %s %s: dispatcher is nil", method, path)
	}
	if path == "" {
		return fmt.Errorf("register http route %s %s: path is required", method, path)
	}
	if action == "" {
		return fmt.Errorf("register http route %s %s: action is required", method, path)
	}

	pattern := buildServeMuxPattern(method, path)
	pathParamsExtractor := serveMuxPathParamsExtractor(extractPathParamNames(path))
	routePatternExtractor := func(req *stdhttp.Request) string {
		return path
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("register http route %s %s: %v", method, path, recovered)
		}
	}()

	r.mux.HandleFunc(pattern, r.dispatcher.handle(action, pathParamsExtractor, routePatternExtractor))
	return nil
}

func buildServeMuxPattern(method string, path string) string {
	if method == "" {
		return path
	}
	return method + " " + path
}

func serveMuxPathParamsExtractor(paramNames []string) PathParamsExtractor {
	if len(paramNames) == 0 {
		return nil
	}

	return func(r *stdhttp.Request) map[string]string {
		params := make(map[string]string, len(paramNames))
		for _, name := range paramNames {
			if value := r.PathValue(name); value != "" {
				params[name] = value
			}
		}
		if len(params) == 0 {
			return nil
		}
		return params
	}
}

func extractPathParamNames(path string) []string {
	if path == "" {
		return nil
	}

	segments := strings.Split(path, "/")
	params := make([]string, 0, len(segments))
	for _, segment := range segments {
		if len(segment) < 3 || segment[0] != '{' || segment[len(segment)-1] != '}' {
			continue
		}

		name := strings.TrimSuffix(strings.TrimPrefix(segment, "{"), "}")
		name = strings.TrimSuffix(name, "...")
		if name == "" || name == "$" {
			continue
		}
		params = append(params, name)
	}

	if len(params) == 0 {
		return nil
	}
	return params
}
