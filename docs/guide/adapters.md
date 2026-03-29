# /adapters
Adapters define how your project connects to the outside world.
This can include inbound adapters like HTTP or WebSockets, but also outbound adapters like external APIs you consume.



## Defining an adapter
Defining an adapter is straightforward.


### 1. Create a init.go with `func Init()`
```go
package myAdapter

func Init(inj injection.Container) error {// [!code ++:3]
  return inj.Register(openai.NewClient(...))
}
```
This is just a minimal example. you may also want to read config, create instances etc.

### 2. Register in main.go
```go
package main

import (
 //...
)

var initializers = []any{
	//------------  Core domains and packages
	core_infra.Init,
    myAdapter.Init, // [!code focus] [!code ++]
	//------------ Domains Infra
	// ----------- Domain Apps
	
	//------------ Initialize APIs
	api.Init,
}

func main() {
	injector := injection.NewInjector()

	for _, init := range initializers {
		err = injector.Invoke(init)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
	}
}
```
::: info
The exact place where you add an adapter may be different. Some adapters need to be initialized before `api.Init`, others after it.
:::

Now dependency injection sees your adapter and initializes it.
