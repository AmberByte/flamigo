## Directory Structure

Flamigo projects follow an opinionated, domain-driven folder structure that promotes clean architecture and separation of concerns. This layout helps developers quickly locate and reason about different parts of the system, especially in larger codebases:

```
cmd/                     # Entry point of the application (main executable)

internal/
  api/                   # API strategy implementations
  
  domains/               # Core business logic organized by domain
    <domainName>/
      app/               # Application services and domain listeners
                         # Listens to events from other domains
                         # (package <domainName>_app)

      domain/            # Domain layer with entities, interfaces, value objects, and events
                         # (package <domainName>)

      infrastructure/    # Infrastructure layer with implementations for repositories, adapters, etc.
                         # (package <domainName>_infra)

  adapters/              # Adapters to the outside world, such as HTTP, WebSockets, or external APIs
```

This structure is designed to support modular development and enforce the boundaries between domain logic and infrastructure, making your codebase both testable and scalable.

## Linking Code Together
The whole framework links everything together by using dependency injection. the pattern is the following:
Different packages define an `init.go` file with a method like the following to handle dependency injection and initialize the package.
### init.go
```go
package api
func Init(inj injection.Container, db database.Database) error {
  // Does all the initialization
}
```

This can then be invoked through the injector:
```go
container := injection.NewInjector()
container.Invoke(api.Init)
```
::: tip
The injector can also inject itself, which allows nesting `Init` functions inside each other.
:::

It is common to have many of these across an app, and they may call each other in a nested way.

e.g. in a small real world app this might be the structure of init functions calling each other:
```
main.go 
        -> api/init.go
                      -> api/users/init.go
                      -> api/messages/init.go
        -> domains/users/infrastructure/init.go
        -> domains/messages/infrastructure/init.go

        -> domains/users/app/init.go
        -> domains/users/messages/init.go

        -> adapters/websocket/init.go
```

::: warning
As already pointed out under Dependency Injection, only use this during startup and never during runtime of your app
:::
