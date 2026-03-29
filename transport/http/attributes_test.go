package http_test

import (
	stdhttp "net/http"
	"net/url"
	"testing"

	flamigo "github.com/amberbyte/flamigo/core"
	transporthttp "github.com/amberbyte/flamigo/transport/http"
	"github.com/amberbyte/flamigo/strategies"
	"github.com/stretchr/testify/assert"
)

func TestAttachRequestMetadata(t *testing.T) {
	ctx := strategies.NewContext(flamigo.NewContext(nil, flamigo.NewServerActor("test")), "app::users:get", nil)

	transporthttp.AttachRequestMetadata(ctx.Request(), transporthttp.Metadata{
		Method:     "GET",
		PathParams: map[string]string{"id": "42"},
		Query:      url.Values{"expand": {"profile"}},
		Headers:    stdhttp.Header{"X-Test": []string{"yes"}},
		Route:      "/users/{id}",
	})

	method, ok := transporthttp.Method(ctx)
	assert.True(t, ok)
	assert.Equal(t, "GET", method)

	id, ok := transporthttp.PathParam(ctx, "id")
	assert.True(t, ok)
	assert.Equal(t, "42", id)

	expand, ok := transporthttp.QueryParam(ctx, "expand")
	assert.True(t, ok)
	assert.Equal(t, "profile", expand)

	header, ok := transporthttp.Header(ctx, "X-Test")
	assert.True(t, ok)
	assert.Equal(t, "yes", header)

	route, ok := transporthttp.Route(ctx)
	assert.True(t, ok)
	assert.Equal(t, "/users/{id}", route)
}
