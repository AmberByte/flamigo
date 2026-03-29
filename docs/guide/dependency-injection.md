# Dependency Injection

Flamigo embraces **Dependency Injection (DI)** as a core design principle to promote loose coupling and enhance testability across your application. Instead of hard-coding dependencies directly inside your services or domains, Flamigo encourages the use of well-defined interfaces and provides a mechanism to inject concrete implementations where needed.

All dependency resolution is handled at runtime, enabling you to easily substitute mocks, switch implementations, or reorganize your architecture without refactoring core business logic. By decoupling components from their dependencies, Flamigo helps you write more maintainable, modular, and scalable code.

---

# Dependency Manager

Flamigo ships with a simple but powerful built-in DI **Container**, located in the `injection` package. It handles wiring and resolving dependencies for your application at runtime.

## Getting Started

To create a new injector instance:

```go
package main

import (
  "github.com/amberbyte/flamigo/injection"
)

func main() {
  inj := injection.NewInjector()
}
```

::: danger
Only add and inject dependencies during your apps startup, never on demand during runtime as this may lead to unforseen problems and hard to track down runtime errors
:::

## Adding a Dependency

You can register a dependency using `Register(i any) error`:

```go
package main

import (
  "github.com/amberbyte/flamigo/injection"
)

func main() {
  inj := injection.NewInjector()
  myDep := &ExampleRepository{}

  if err := inj.Register(myDep); err != nil {
    panic(err)
  }
}
```

:::info
An error is returned if a dependency is already managed or conflicts with another injectable.
:::

## Injecting Dependencies

Dependencies can be dynamically injected by defining a function that takes them as parameters:

```go
func createMyLogic(repo *ExampleRepository) error {
  // your logic here
  return nil
}
```

Then invoke it through the injector:

```go
func main() {
  inj := injection.NewInjector()
  inj.Register(&ExampleRepository{})

  if err := inj.Invoke(createMyLogic); err != nil {
    panic(err)
  }
}
```

:::info
Interfaces are also supported — the Dependency Manager automatically matches the appropriate implementation if it satisfies the interface.
:::

### `Invoke` returns an error if:

- A required dependency cannot be resolved  
- The function returns an error  
- There’s ambiguity — i.e., multiple injectables satisfy the same interface  

Call-time `Invoke(..., args...)` arguments take precedence over registered injectables. That lets startup code override dependencies explicitly for a single execution without mutating the container.
