package http_test

import (
	"context"
	"errors"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/amberbyte/flamigo/strategies"
	transporthttp "github.com/amberbyte/flamigo/transport/http"
	"github.com/stretchr/testify/assert"
)

type stubRouter struct {
	invoke func(ctx strategies.Context) strategies.Result
}

func (s *stubRouter) Register(namespace string, registry strategies.AppRegistry) error {
	return nil
}

func (s *stubRouter) Invoke(ctx strategies.Context) strategies.Result {
	return s.invoke(ctx)
}

func TestDispatcherHandleSuccess(t *testing.T) {
	router := &stubRouter{
		invoke: func(ctx strategies.Context) strategies.Result {
			method, ok := transporthttp.Method(ctx)
			assert.True(t, ok)
			assert.Equal(t, stdhttp.MethodGet, method)

			id, ok := transporthttp.PathParam(ctx, "id")
			assert.True(t, ok)
			assert.Equal(t, "42", id)

			ctx.Response().SetPayload(map[string]string{"ok": "true"})
			return ctx.Response()
		},
	}

	dispatcher := transporthttp.NewDispatcher(
		router,
		transporthttp.WithPathParamsExtractor(func(r *stdhttp.Request) map[string]string {
			return map[string]string{"id": "42"}
		}),
	)

	req := httptest.NewRequest(stdhttp.MethodGet, "/users/42", nil)
	w := httptest.NewRecorder()
	dispatcher.Handle("app::users:get")(w, req)

	assert.Equal(t, stdhttp.StatusOK, w.Code)
	assert.JSONEq(t, `{"ok":"true"}`, w.Body.String())
}

func TestDispatcherHandlePublicError(t *testing.T) {
	router := &stubRouter{
		invoke: func(ctx strategies.Context) strategies.Result {
			ctx.Response().SetError(flamigo.NewError("not allowed", flamigo.Public("Forbidden"), flamigo.StatusCode(stdhttp.StatusForbidden)))
			return ctx.Response()
		},
	}

	dispatcher := transporthttp.NewDispatcher(router)
	req := httptest.NewRequest(stdhttp.MethodGet, "/users/42", nil)
	w := httptest.NewRecorder()
	dispatcher.Handle("app::users:get")(w, req)

	assert.Equal(t, stdhttp.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Forbidden")
}

func TestDispatcherActorFactoryError(t *testing.T) {
	router := &stubRouter{
		invoke: func(ctx strategies.Context) strategies.Result {
			t.Fatal("router should not be called")
			return nil
		},
	}

	dispatcher := transporthttp.NewDispatcher(router, transporthttp.WithActorFactory(func(r *stdhttp.Request, w stdhttp.ResponseWriter) (flamigo.Actor, error) {
		return nil, errors.New("boom")
	}))
	req := httptest.NewRequest(stdhttp.MethodGet, "/users/42", nil)
	w := httptest.NewRecorder()
	dispatcher.Handle("app::users:get")(w, req)

	assert.Equal(t, stdhttp.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "boom")
}

func TestDefaultActorFactoryUsesBackendActor(t *testing.T) {
	router := &stubRouter{
		invoke: func(ctx strategies.Context) strategies.Result {
			assert.Equal(t, flamigo.TypeActorServer, ctx.Actor().Type())
			assert.Equal(t, context.Background().Err(), ctx.Err())
			return ctx.Response()
		},
	}

	dispatcher := transporthttp.NewDispatcher(router)
	req := httptest.NewRequest(stdhttp.MethodGet, "/users/42", nil)
	w := httptest.NewRecorder()
	dispatcher.Handle("app::users:get")(w, req)

	assert.Equal(t, stdhttp.StatusNoContent, w.Code)
}
