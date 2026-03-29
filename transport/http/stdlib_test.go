package http_test

import (
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/amberbyte/flamigo/strategies"
	transporthttp "github.com/amberbyte/flamigo/transport/http"
	"github.com/stretchr/testify/assert"
)

func TestServeMuxRegistrarHandle(t *testing.T) {
	router := &stubRouter{
		invoke: func(ctx strategies.Context) strategies.Result {
			method, ok := transporthttp.Method(ctx)
			if !assert.True(t, ok) {
				return ctx.Response()
			}
			assert.Equal(t, stdhttp.MethodGet, method)

			id, ok := transporthttp.PathParam(ctx, "id")
			if !assert.True(t, ok) {
				return ctx.Response()
			}
			assert.Equal(t, "42", id)

			route, ok := transporthttp.Route(ctx)
			if !assert.True(t, ok) {
				return ctx.Response()
			}
			assert.Equal(t, "/users/{id}", route)

			ctx.Response().SetPayload(map[string]string{"id": id})
			return ctx.Response()
		},
	}

	dispatcher := transporthttp.NewDispatcher(router)
	mux := stdhttp.NewServeMux()
	registrar := transporthttp.NewServeMuxRegistrar(mux, dispatcher)

	err := registrar.Handle(stdhttp.MethodGet, "/users/{id}", "app::users:get")
	if !assert.NoError(t, err) {
		return
	}

	req := httptest.NewRequest(stdhttp.MethodGet, "/users/42", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, stdhttp.StatusOK, w.Code)
	assert.JSONEq(t, `{"id":"42"}`, w.Body.String())
}

func TestServeMuxRegistrarRejectsNilDependencies(t *testing.T) {
	dispatcher := transporthttp.NewDispatcher(&stubRouter{
		invoke: func(ctx strategies.Context) strategies.Result {
			return ctx.Response()
		},
	})

	err := transporthttp.NewServeMuxRegistrar(nil, dispatcher).Handle(stdhttp.MethodGet, "/users/{id}", "app::users:get")
	assert.Error(t, err)

	err = transporthttp.NewServeMuxRegistrar(stdhttp.NewServeMux(), nil).Handle(stdhttp.MethodGet, "/users/{id}", "app::users:get")
	assert.Error(t, err)
}
