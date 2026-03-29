package http

import (
	stdhttp "net/http"
	"net/url"

	"github.com/amberbyte/flamigo/strategies"
)

const (
	AttrMethod     = "http.method"
	AttrPathParams = "http.path_params"
	AttrQuery      = "http.query"
	AttrHeaders    = "http.headers"
	AttrRoute      = "http.route"
)

type Metadata struct {
	Method     string
	PathParams map[string]string
	Query      url.Values
	Headers    stdhttp.Header
	Route      string
}

func AttachRequestMetadata(req *strategies.Request, metadata Metadata) {
	if metadata.Method != "" {
		req.SetAttribute(AttrMethod, metadata.Method)
	}
	if metadata.PathParams != nil {
		req.SetAttribute(AttrPathParams, metadata.PathParams)
	}
	if metadata.Query != nil {
		req.SetAttribute(AttrQuery, metadata.Query)
	}
	if metadata.Headers != nil {
		req.SetAttribute(AttrHeaders, metadata.Headers)
	}
	if metadata.Route != "" {
		req.SetAttribute(AttrRoute, metadata.Route)
	}
}

func Method(ctx strategies.Context) (string, bool) {
	value, ok := ctx.Request().Attribute(AttrMethod)
	if !ok {
		return "", false
	}
	method, ok := value.(string)
	return method, ok
}

func PathParams(ctx strategies.Context) (map[string]string, bool) {
	value, ok := ctx.Request().Attribute(AttrPathParams)
	if !ok {
		return nil, false
	}
	params, ok := value.(map[string]string)
	return params, ok
}

func PathParam(ctx strategies.Context, key string) (string, bool) {
	params, ok := PathParams(ctx)
	if !ok {
		return "", false
	}
	value, ok := params[key]
	return value, ok
}

func Query(ctx strategies.Context) (url.Values, bool) {
	value, ok := ctx.Request().Attribute(AttrQuery)
	if !ok {
		return nil, false
	}
	query, ok := value.(url.Values)
	return query, ok
}

func QueryParam(ctx strategies.Context, key string) (string, bool) {
	query, ok := Query(ctx)
	if !ok {
		return "", false
	}
	values, ok := query[key]
	if !ok || len(values) == 0 {
		return "", false
	}
	return values[0], true
}

func Headers(ctx strategies.Context) (stdhttp.Header, bool) {
	value, ok := ctx.Request().Attribute(AttrHeaders)
	if !ok {
		return nil, false
	}
	headers, ok := value.(stdhttp.Header)
	return headers, ok
}

func Header(ctx strategies.Context, key string) (string, bool) {
	headers, ok := Headers(ctx)
	if !ok {
		return "", false
	}
	value := headers.Get(key)
	if value == "" {
		return "", false
	}
	return value, true
}

func Route(ctx strategies.Context) (string, bool) {
	value, ok := ctx.Request().Attribute(AttrRoute)
	if !ok {
		return "", false
	}
	route, ok := value.(string)
	return route, ok
}
