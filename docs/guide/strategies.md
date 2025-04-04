# Strategies

Flamigo includes a built-in **strategy resolution system** that helps you decouple and organize your API layer logic in a clean, modular, and reusable way.

A **strategy** in Flamigo is a named, callable function that encapsulates a specific piece of logic — typically used to handle external inputs such as API calls, WebSocket messages, or other interaction layers. Rather than binding behavior directly to a transport (like HTTP routes), Flamigo lets you define strategies under unique names and resolve them dynamically at runtime.

This design offers:

- A unified abstraction over different transport layers (HTTP, WebSocket, CLI, etc.)
- The ability to call strategies from anywhere — including other strategies or domain listeners
- A testable, reusable structure for external interfaces
- A generalized pattern that can be reused across different domains and contexts

---

## Strategy Signatures

A strategy is simply a function that accepts a `strategies.Context` — an extended form of Flamigo’s context that includes additional request and response utilities.

```go
type StrategyFunc func(ctx strategies.Context)
```

---

## Registering Strategies

Strategies are registered by name — usually in the `internal/api/` layer — and grouped logically by application area or domain.

```go
registry := strategies.NewRegistry("app")

registry.Register("app::domain:doSth", func(ctx strategies.Context) {
  // Strategy logic here
})
```

You can register multiple strategies under one registry, following a consistent naming convention (e.g. `app::domain:action`).

---

## Strategy Context

The strategy context extends the base `flamigo.Context`, so it includes access to the current `Actor`. In addition, it provides structured access to the incoming request and the response.

```go
type Request interface {
  Action() string      // The action name (e.g. app::domain:doSth)
  Payload() any        // The raw payload sent to the strategy
  Bind(target any) error // Decode the payload (typically JSON) into a target struct
}

type Response interface {
  SetResult(payload any) string // Sets the result to return from the strategy
  SetError(err error)           // Sets an error as the result
  Result() any
  Err() error
  IsOk() bool
  IsError() bool
}
```

This setup provides everything your strategy needs to process input, work with actors, and respond in a consistent and structured way.

---

## Calling Strategies

To call a strategy, you compose a strategy context with the appropriate `flamingo.Context`, action and payload, then invoke it through the registry:

```go
registry := strategies.NewRegistry("app")

ctx := strategies.NewContext(flamigoCtx, "app::foo:bar", `{"foo": "bar"}`)

result := registry.Use(ctx)

if result.IsOk() {
  // Handle success
}
// ...
```

This approach enables strategies to be invoked internally — from domains, listeners, or even other strategies — without tightly coupling to the transport layer or request format.