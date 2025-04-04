# /interfaces
interfaces defines interfaces your project has to the outside world.
This could be other apis you consume (e.g. openai etc.) but this also has external endpoints you provide (e.g. websockets)



## Defining a interface
Defining a interface is straightforward


### 1. Create a init.go with `func Init()`
```go
package myInterface

func Init(inj injection.DependencyManager) error {// [!code ++:3]
  return inj.AddInjectable(openai.NewClient(...))
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
    myInterface.Init, // [!code focus] [!code ++]
	//------------ Domains Infra
	// ----------- Domain Apps
	
	//------------ Initialize APIs
	api.Init,
	websocket.Init,
}

func main() {
	injector := injection.NewInjecter()

	for _, init := range initializers {
		err = injector.Execute(init)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
	}
}
```
::: info
The exact place where you may add your interface may be different (e.g. it may be required at the beginning or at the end)
:::

Now dependency injection sees your interface and initalizes it