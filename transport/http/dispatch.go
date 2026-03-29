package http

import (
	"encoding/json"
	"io"
	stdhttp "net/http"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/amberbyte/flamigo/strategies"
)

type ActorFactory func(r *stdhttp.Request, w stdhttp.ResponseWriter) (flamigo.Actor, error)
type PayloadDecoder func(r *stdhttp.Request) (any, error)
type PathParamsExtractor func(r *stdhttp.Request) map[string]string
type RoutePatternExtractor func(r *stdhttp.Request) string

type DispatchOption func(*Dispatcher)

type Dispatcher struct {
	router               strategies.AppRouter
	actorFactory         ActorFactory
	payloadDecoder       PayloadDecoder
	pathParamsExtractor  PathParamsExtractor
	routePatternExtract  RoutePatternExtractor
	errorEncoder         ErrorEncoder
	payloadEncoder       PayloadEncoder
}

func WithActorFactory(factory ActorFactory) DispatchOption {
	return func(d *Dispatcher) {
		d.actorFactory = factory
	}
}

func WithPayloadDecoder(decoder PayloadDecoder) DispatchOption {
	return func(d *Dispatcher) {
		d.payloadDecoder = decoder
	}
}

func WithPathParamsExtractor(extractor PathParamsExtractor) DispatchOption {
	return func(d *Dispatcher) {
		d.pathParamsExtractor = extractor
	}
}

func WithRoutePatternExtractor(extractor RoutePatternExtractor) DispatchOption {
	return func(d *Dispatcher) {
		d.routePatternExtract = extractor
	}
}

func WithErrorEncoder(encoder ErrorEncoder) DispatchOption {
	return func(d *Dispatcher) {
		d.errorEncoder = encoder
	}
}

func WithPayloadEncoder(encoder PayloadEncoder) DispatchOption {
	return func(d *Dispatcher) {
		d.payloadEncoder = encoder
	}
}

func defaultActorFactory(r *stdhttp.Request, w stdhttp.ResponseWriter) (flamigo.Actor, error) {
	return flamigo.NewServerActor("http"), nil
}

func defaultPayloadDecoder(r *stdhttp.Request) (any, error) {
	if r.Body == nil {
		return nil, nil
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return json.RawMessage(data), nil
}

func NewDispatcher(router strategies.AppRouter, opts ...DispatchOption) *Dispatcher {
	dispatcher := &Dispatcher{
		router:         router,
		actorFactory:   defaultActorFactory,
		payloadDecoder: defaultPayloadDecoder,
		errorEncoder:   DefaultErrorEncoder,
		payloadEncoder: DefaultPayloadEncoder,
	}
	for _, opt := range opts {
		opt(dispatcher)
	}
	return dispatcher
}

func (d *Dispatcher) Handle(action string) stdhttp.HandlerFunc {
	return d.handle(action, d.pathParamsExtractor, d.routePatternExtract)
}

func (d *Dispatcher) handle(action string, pathParamsExtractor PathParamsExtractor, routePatternExtractor RoutePatternExtractor) stdhttp.HandlerFunc {
	return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		actor, err := d.actorFactory(r, w)
		if err != nil {
			d.errorEncoder(w, r, err)
			return
		}

		payload, err := d.payloadDecoder(r)
		if err != nil {
			d.errorEncoder(w, r, err)
			return
		}

		appCtx := flamigo.NewContext(r.Context(), actor)
		strategyCtx := strategies.NewContext(appCtx, action, payload)
		AttachRequestMetadata(strategyCtx.Request(), Metadata{
			Method:     r.Method,
			PathParams: extractPathParams(pathParamsExtractor, r),
			Query:      r.URL.Query(),
			Headers:    r.Header.Clone(),
			Route:      extractRoutePattern(routePatternExtractor, r),
		})

		result := d.router.Invoke(strategyCtx)
		if err := result.Err(); err != nil {
			d.errorEncoder(w, r, err)
			return
		}

		d.payloadEncoder(w, r, result.Payload())
	}
}

func extractPathParams(extractor PathParamsExtractor, r *stdhttp.Request) map[string]string {
	if extractor == nil {
		return nil
	}
	return extractor(r)
}

func extractRoutePattern(extractor RoutePatternExtractor, r *stdhttp.Request) string {
	if extractor == nil {
		return ""
	}
	return extractor(r)
}
