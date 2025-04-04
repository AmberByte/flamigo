# /domains
Domains is where your business logic resides.
A core concept of flamigo is to structure your business logic around business domains and not tightly coupling them.
This is where domains come into play.

## Folder structure
Every domain lives under its own directory under `/domains`.
A domain named `messages` therefore lives under `/domains/messages`

Every domain is then composed of three sub directories:

- **`domain`**  
  This folder contains interfaces and contracts your domain has to the outside world. this can be entities, events, repository interfaces etc. it normally contains a package named after the domain.  
  In our messages example it would be `package messages`
- **`app`**  
  This contains event listeners on other domains and app logic the domain contains.  
  In our messages example this would be the package `package messages_app`
- **`infrastructure`**  
  This contains implementation for repositories and binds the domain together.  
  Package would be  `package messages_infra`

## Domain Aggregate
To simplify dependency injection between domains they may define a aggregate in the domain folder that contains all injectables (most of the time repositories).
### 1. Define the interface in the domain
```go
package messages

type Domain interface {// [!code ++:4]
  Messages() MessageRepository
  Addresses() AddressRepository
}
```
### 2. Implement domain aggregate `messages_infra`
As you see above we only defined an interface for now. in messages_infra we now bind it together in infrastructure directory

```go
package messages_infra

var _ messages.Domain = (*domain)(nil)// [!code ++:14]

type domain struct {
  msg messages.MessageRepository
  adr messages.AddressRepository
}

func (d *domain) Messages() messages.MessageRepository {
  return msg
}

func (d *domain) Addresses() messages.AddressRepository {
  return msg
}
```

### 3. Now add a `func Init()` to be able to use it in dependency injection
```go
package messages_infra

func Init(inj injection.DependencyManager, db database.Database) error {// [!code ++:11]
  messageRepo := newMessageRepo(db)
  addressRepo := newAddressRepo(db)

  dmn := &domain{
    msg: messageRepo
    adr: addressRepo
  }

  return inj.AddInjectable(dmn)
}

var _ messages.Domain = (*domain)(nil)

type domain struct {
  msg messages.MessageRepository
  adr messages.AddressRepository
}

func (d *domain) Messages() messages.MessageRepository {
  return msg
}

func (d *domain) Addresses() messages.AddressRepository {
  return msg
}
```

### 4. Now you can use this in `cmd/main.go` to inject

```go
package main

import (
 //...
)

var initializers = []any{
	//------------  Core domains and packages
	core_infra.Init,
	//------------ Domains Infra
	// ----------- Domain Apps
  messages_infra.Init, // [!code focus] [!code ++]
	
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
now the message_infra init is run during startup phase, and creates repositories and injects the domain aggregate