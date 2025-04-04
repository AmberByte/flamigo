# /api
API defines how others interact with your domain.
The API can be structured around your domains or follow any other approach, theres no limitation you can decided what suits you best.

## Defining strategies in api layer
Defining a strategy in api layer can be done in a special dependency injection wrapper.
Its convenient to do everything in there, but for larger and more complex stragtegies it may be reasonable to split it up.

### 1. Define strategy
```go
func createStrategyGetMessages(strategy strategy.Registry, msgDomain messages.Domain) {// [!code ++:5]
  strategy := func(ctx strategy.Context) {
    // Do some logic here
  }
}
```

### 2. Register your strategy in the registry
```go
func createStrategyGetMessages(strategy strategy.Registry, msgDomain messages.Domain) {
  strategy := func(ctx strategy.Context) {
    // Do some logic here
  }

  strategy.Register("app::messages:get", strategy)// [!code ++]
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

func Init(inj injection.DependencyManager) error {
	return inj.ExecuteList(apiModules)
}

```
Now your api is registered

## Limiting Based on the Actor
You can also limit based on the actor thas coming in. this can be extended with your own validators to your needs:
```go
func createStrategyGetMessages(strategy strategy.Registry, msgDomain messages.Domain) {
  strategy := func(ctx strategy.Context) {// [!code focus:9]
    err := flamigo.RequireActorWithClaims[flamigo.Actor](ctx, flamigo.IsServer())// [!code ++:5]
    if err != nil {
      ctx.Response.SetError(err)
      return
    }
    // Do some logic here
  }

  strategy.Register("app::messages:get", strategy)
}