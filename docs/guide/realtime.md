# Realtime

Flamigo comes with built-in support for **realtime event handling**, allowing your application to react instantly to domain events and push updates to the frontend as they happen.

At the core of this system is the **Realtime Event Bus**, which serves as a central communication channel for **domain events**. External interfaces like WebSockets can subscribe to those events and project them into UI notifications, live updates, or other transport-specific messages.

By leveraging the event bus, Flamigo keeps internal domain reactions decoupled while still making it straightforward to build realtime client updates on top.

# The Event Bus
All logic for sending realtime events is in `realtime` package

```go
package main

import (
  "github.com/amberbyte/flamigo/realtime"
)

func main() {
  bus := realtime.NewBus()
}
```

# Events
flamigo has its own event type:
```go
type Event interface {
  Topics() []Topic
}
```

A event must carry a list of topics (at least one).

## Topics
A Topic is composed hierachically (e.g. users/userId/xuah4z47fs)

```go
realtime.NewTopic("users", "userId", "xuah4z47fs")
```

## Publishing events
You can publish events by calling `Publish` method.
Publishing events is asynchronous with respect to the listener work, but the bus applies backpressure when subscriber queues fill instead of dropping events silently.
If you need to wait until all matching listeners have processed an event, use `PublishSync`.

```go
bus.Publish(myevent)
```

## Recieving events
You can receive on topics by creating a listener function and then subscribing to topics

```go
listener := func(ctx realtime.Context, evt realtime.Topci) {
  ...
}
subscription := bus.Subscribe(listener)
subscription.SubscripeTopic("foo/bar")
```

## Wildcard subscriptions
Topic matching is exact by default.
Events support wildcards with `*` per path segment when you want to subscribe more broadly.
```go
listener := func(ctx realtime.Context, evt realtime.Topci) {
  ...
}
subscription := bus.Subscribe(listener)
subscription.SubscripeTopic("foo/*") // Subscribes to foo/<single-segment>
```
For example this is useful when you want to subscribe to all user id topics and not a specific user id: `userId/*`

## Client Messages
Domain events can be projected to client messages in your interface layer.
Flamigo still provides a small transport helper in `github.com/amberbyte/flamigo/realtime/client`:
```go
type Message interface {
	Topic() string
	Payload() any
	MarshalClientPayload() ([]byte, error)
}
```
In practice, a WebSocket client should subscribe to domain events and translate selected events into `client.Message`s before sending them to the frontend.
