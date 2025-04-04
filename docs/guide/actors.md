# Actors

In Flamigo, **everything happens through an Actor**. Actors are a fundamental concept in the framework — they represent *who* is executing a piece of logic and can carry additional context such as authentication data, session metadata, or the type of connection.

This design provides a clean and unified way to propagate identity and context throughout your application. Whether you're building APIs, event listeners, or realtime features, knowing *who* triggered the action helps your logic adapt accordingly.

---

## Standard Actors

Flamigo includes two built-in actor types:

- **Backend Actor**  
  The default actor representing the system itself — used for internal processes like scheduled jobs, backend-triggered events, or inter-domain communication.

- **WebSocket Actor**  
  Represents a live WebSocket connection, typically tied to a user or session. This actor enables real-time interactions between the frontend and backend while retaining contextual awareness of the client.

You can also define **custom actor types** to model different execution identities — such as authenticated users, API tokens, automation tools, or external systems.

---

## Retrieving an Actor

Flamigo extends Go's standard `context.Context` with additional functionality to retrieve the current actor. This extended context can be used anywhere a regular `context.Context` is expected.

```go
type Context interface {
  context.Context
  Actor() Actor
}
```

By using this extended context, you gain seamless access to actor-related metadata in your services, listeners, or handlers.

Theres helpers to check if a actor exists:
```go
func RequireActorWithClaims[T Actor](ctx Context, opts ...ActorClaimValidator) (parsedActor T, err error)
```
---

## Actor Interfaces

### `Actor` Interface

The base interface all actors must implement:

```go
type Actor interface {
  Type() string
}
```

The `Type()` method identifies the actor type (e.g., `"backend"`, `"websocket"`, `"user"`), allowing you to branch logic accordingly.

---

### `UserActor` Interface (Auth Plugin)

If you're using Flamigo’s authentication plugin, you can work with a specialized `UserActor` interface that exposes the authenticated user:

```go
type UserActor interface {
  Actor
  User() *auth.User
}
```

This allows your services to securely access the current user’s identity and permissions when needed.
