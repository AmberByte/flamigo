# HTTP

Flamigo provides reusable HTTP transport primitives in `transport/http` and can scaffold an HTTP adapter for generated projects under `internal/adapters/http`.

The goal is to keep transport mechanics in the framework while keeping route decisions, actor policy, authentication, and payload shaping in your app.

::: warning
The generated HTTP adapter is optional and only exists if you enable the `transport_http` feature during project initialization.
:::

## What the framework provides

The framework HTTP package provides:

- a strategy dispatcher for `method + path + action` flow
- request metadata attachment for strategy contexts
- default payload and error encoding
- a reusable `Registrar` abstraction
- a standard library `ServeMux` registrar

At the framework level, HTTP requests are translated into:

- an `Actor`
- a strategy action such as `app::messages:get`
- a request payload
- request attributes like:
  - `http.method`
  - `http.path_params`
  - `http.query`
  - `http.headers`
  - `http.route`

## What the generated app adapter provides

When enabled in the project template, `internal/adapters/http` creates a small app-owned server wrapper around `transport/http`.

That adapter is responsible for:

- creating the concrete HTTP actor for your app
- holding the router/server instance
- exposing a route registrar for startup wiring
- starting the HTTP server

## Registering routes next to strategies

The intended pattern is to register routes close to the strategy that handles them:

```go
func createStrategyGetMessages(
  registry strategies.Registry[strategies.Context],
  routes http.Registrar,
  msgDomain messages.Domain,
) error {
  handler := func(ctx strategies.Context) {
    // strategy logic
  }

  if err := registry.Register("messages:get", handler); err != nil {
    return err
  }

  return routes.Handle("GET", "/messages/{id}", "app::messages:get")
}
```

This keeps route declarations close to the application action without moving HTTP concerns into the `strategies` package.

## Accessing HTTP request metadata

HTTP-specific details are exposed through request attributes on the strategy context, not through transport-specific strategy methods.

Examples:

```go
method, _ := transporthttp.Method(ctx)
id, _ := transporthttp.PathParam(ctx, "id")
route, _ := transporthttp.Route(ctx)
```

This keeps the strategy API transport-agnostic while still allowing adapters to provide useful metadata.

## Payload binding and validation

The framework binds raw request payloads into the strategy request, but validation is left to the app.

In generated projects, the template provides an app-owned helper in `pkg/validation`:

```go
if err := validation.BindAndValidate(ctx.Request(), &payload); err != nil {
  ctx.Response().SetError(err)
  return
}
```

That keeps validation policy out of Flamigo core while still providing a convenient default in the scaffolded app.

## What stays app-owned

The framework does not decide:

- your authentication flow
- your validation rules
- your DTO design
- your route organization
- your error policy beyond basic transport helpers

Those remain adapter- and app-level concerns.
