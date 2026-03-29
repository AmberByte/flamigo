# /api
API defines how others interact with your domain.
The API can be structured around your domains or follow any other approach, theres no limitation you can decided what suits you best.

## Defining strategies in api layer
Defining a strategy in api layer can be done in a special dependency injection wrapper.
Its convenient to do everything in there, but for larger and more complex stragtegies it may be reasonable to split it up.

### 1. Define strategy
```go
func createStrategyGetMessages(registry strategies.Registry[strategies.Context], msgDomain messages.Domain) {// [!code ++:5]
  handler := func(ctx strategy.Context) {
    // Do some logic here
  }
}
```

### 2. Register your strategy in the registry
```go
func createStrategyGetMessages(registry strategies.Registry[strategies.Context], msgDomain messages.Domain) error {
  handler := func(ctx strategy.Context) {
    // Do some logic here
  }

  return registry.Register("messages:get", handler)// [!code ++]
}
```

### 3. Register your strategy in api
You can now add the method to your api func Init()
```go
package api

import (
	"github.com/amberbyte/flamigo/injection"
)

var apiModules = []any{
  createStrategyGetMessages,// [!code ++]
}

func Init(inj injection.Container) error {
	return inj.InvokeAll(apiModules)
}

```
Now your api is registered

### Optional: Register HTTP Routes Next To Strategies

If a strategy is exposed over HTTP, inject an HTTP route registrar from your HTTP interface layer and bind the route next to the strategy registration:

```go
func createStrategyGetMessages(
  registry strategies.Registry[strategies.Context],
  routes httptransport.Registrar,
  msgDomain messages.Domain,
) error {
  handler := func(ctx strategies.Context) {
    // Do some logic here
  }

  if err := registry.Register("messages:get", handler); err != nil {
    return err
  }

  return routes.Handle("GET", "/messages/{id}", "app::messages:get")
}
```

This keeps HTTP method/path mapping close to the strategy without moving HTTP concerns into the `strategies` package itself.

## Limiting Based on the Actor
You can also limit based on the actor thas coming in. this can be extended with your own validators to your needs:
```go
func createStrategyGetMessages(registry strategies.Registry[strategies.Context], msgDomain messages.Domain) error {
  handler := func(ctx strategy.Context) {// [!code focus:9]
    err := flamigo.RequireActorWithClaims[flamigo.Actor](ctx, flamigo.IsServer())// [!code ++:5]
    if err != nil {
      ctx.Response.SetError(err)
      return
    }
    // Do some logic here
  }

  return registry.Register("messages:get", handler)
}
```
