# Strategies

Flamigo includes a built-in **strategy resolution system** that helps you decouple and organize your API layer logic in a clean, modular, and reusable way.

A **strategy** in Flamigo is a named, callable function that encapsulates a specific piece of logic — typically used to handle external inputs such as API calls, WebSocket messages, or other interaction layers. Rather than binding behavior directly to a transport (like HTTP routes), Flamigo lets you define strategies under unique names and resolve them dynamically at runtime.

This design offers:

- A unified abstraction over different transport layers (HTTP, WebSocket, CLI, etc.)
- The ability to call strategies from anywhere — including other strategies or domain listeners
- A testable, reusable structure for external adapters
- A generalized pattern that can be reused across different domains and contexts

---

## Strategy Signatures

A strategy is simply a function that accepts a `strategies.Context` — an extended form of Flamigo’s context that includes additional request and response utilities.

```go
type StrategyFunc func(ctx strategies.Context)
```

---

## Registering Strategies

Strategies are registered by local action name — usually in the `internal/api/` layer — and grouped logically by application area or domain.

```go
appRegistry := strategies.NewRegistry[strategies.Context]()

if err := appRegistry.Register("domain:doSth", func(ctx strategies.Context) {
  // Strategy logic here
}); err != nil {
  panic(err)
}
```

`Register` expects a local action name such as `domain:action`. It returns an error for invalid names or duplicate registrations.

::: warning
Register registries and strategies during startup only. The strategies package is configured to be wired up at application initialization time, not mutated during normal runtime.
:::

If you want to manage multiple namespaces, use a router above the local registries:

```go
appRegistry := strategies.NewRegistry[strategies.Context]()
adminRegistry := strategies.NewRegistry[strategies.Context]()

router := strategies.NewRouter[strategies.Context]()
router.Register("app", appRegistry)
router.Register("admin", adminRegistry)
```

If you expose strategies over HTTP, keep route binding in the HTTP adapter layer, but register routes close to the strategy itself through an injected route registrar.

---

## Strategy Context

The strategy context extends the base `flamigo.Context`, so it includes access to the current `Actor`. In addition, it provides structured access to the incoming request and the response.

```go
type Request interface {
  Action() string      // The action name (e.g. app::domain:doSth)
  Payload() any        // The raw payload sent to the strategy
  Bind(target any) error // Decode the payload (typically JSON) into a target struct
  SetAttribute(key string, value any)
  Attribute(key string) (any, bool)
}

type Response interface {
  SetPayload(payload any)       // Sets the payload to return from the strategy
  SetError(err error)           // Sets an error as the result
  Payload() any
  Err() error
  IsOk() bool
  IsError() bool
}
```

This setup provides everything your strategy needs to process input, work with actors, and respond in a consistent and structured way.

Transport-specific metadata should be attached as request attributes rather than added directly to the strategy API. For example, an HTTP adapter can attach method, path params, query params, or headers under keys like `http.method`, `http.path_params`, and `http.query`.

---

## Calling Strategies

To call a strategy, you compose a strategy context with the appropriate `flamingo.Context`, full action, and payload, then invoke it through the router:

```go
appRegistry := strategies.NewRegistry[strategies.Context]()
router := strategies.NewRouter[strategies.Context]()
router.Register("app", appRegistry)

ctx := strategies.NewContext(flamigoCtx, "app::foo:bar", `{"foo": "bar"}`)

result := router.Invoke(ctx)

if result.IsOk() {
  // Handle success
}
// ...
```

This approach enables strategies to be invoked internally — from domains, listeners, or even other strategies — without tightly coupling to the transport layer or request format.
