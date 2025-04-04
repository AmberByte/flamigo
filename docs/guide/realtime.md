# Realtime

Flamigo comes with built-in support for **realtime event handling**, allowing your application to react instantly to domain events and push updates to the frontend as they happen.

At the core of this system is the **Realtime Event Bus**, which serves as a central communication channel between domains and external interfaces like WebSockets. This makes it easy to decouple your logic while enabling features like live notifications, dynamic UI updates, or collaborative state changes — all without additional setup.

By leveraging the event bus, Flamigo ensures that both internal domain reactions and external user interfaces stay in sync in real time.

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
Publishing events is asynchronous, and flamigo does not promise anything regarding the order of the events.
Publish also directly returns and does not wait for the events to be completely processed. to do so you can use `PublishSync`. However theres a danger of deadlocking which is why using Publish is prefeered

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
Events support wildcards with `*` for suscribing to everything under this hierachy
```go
listener := func(ctx realtime.Context, evt realtime.Topci) {
  ...
}
subscription := bus.Subscribe(listener)
subscription.SubscripeTopic("foo/*") // Susbcribes to everything under foo
```
For example this is useful when you want to subscribe to all user id topics and not a specific user id: `userId/*`

## Client Messages
Events can be forwarded to clients (e.g. websocket connections) as well.
For this theres the following event interface
```go
type ClientEvent interface {
	ClientMessage() ClientMessage
}
```

A event can fullfil this interface to make it forwardable to clients