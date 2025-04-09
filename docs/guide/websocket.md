# WebSocket (Realtime)

Flamigo's realtime interface is built on **WebSockets**, allowing you to send or push updates from your backend to connected frontend clients — instantly and efficiently.

It integrates directly with the **Realtime Event Bus**, enabling seamless reaction to domain events and bidirectional communication with frontend applications.

> ⚠️ The WebSocket interface must be enabled when initializing a new project.

---

## WebSocket Actor

In line with Flamigo’s actor-based architecture, each WebSocket connection is treated as an **actor**. This allows strategies and domain logic to understand and respond to *who* initiated a given action via WebSocket.

The WebSocket interface provides its own actor type to represent the client.

:::info  
If the **auth** plugin is also enabled, the WebSocket actor will automatically implement the `UserActor` interface — giving you access to user-specific logic out of the box.  
:::

---

### Actor Claim Validators

To help you control behavior based on the type of actor, two claim validators are available:

- **`IsWebsocket`** – Passes only if the actor is a WebSocket actor.  
- **`IsNotWebsocket`** – Fails if the actor is a WebSocket actor.

These can be used in strategies or other actor-aware logic to adapt behavior to the connection type.

---

## Authentication

If the authentication plugin is enabled, the WebSocket interface registers a dedicated **authentication strategy**:

```txt
app::websocket:auth
```

This strategy is triggered when a client connects and attempts to authenticate. You are expected to **customize this strategy** to match your authentication flow — e.g., validate session tokens, headers, or credentials from the client.

---

## The WebSocket Client

Each WebSocket connection is represented by a **WebSocket Client** object. This client:

- Manages the connection’s state
- Tracks the associated user (if authenticated)
- Subscribes to topics on the Realtime Event Bus

By default, client setup is minimal — you are expected to **manually subscribe** clients to topics and customize the client logic as needed.

> 💡 See the `client.go` file for a `TODO` section where you can hook in custom behavior.

---

## Events

The WebSocket interface emits lifecycle events, including:

- **`EventDisconnected`** – Fired when a client disconnects from the server.  
  You can use this to clean up state, revoke sessions, or notify others.

---

## WebSocket Communication

Communication between frontend and backend happens through structured WebSocket messages. Each message follows a standard format:

### Sending a Request

```json
{
  "topic": "app::strategy",
  "payload": {}
}
```

- `topic`: The strategy name to execute on the server.
- `payload`: Arbitrary data to send as input.

---

### Receiving Responses

To receive a response to a specific request, include an `ackId` in the message:

```json
{
  "topic": "app::strategy",
  "payload": {},
  "ackId": "12345"
}
```

The response will be returned with the same `ackId`:

```json
{
  "topic": "app::strategy",
  "payload": {
    "foo": "bar"
  },
  "ackId": "12345"
}
```

This lets the client **match the response to the original request**, which is especially useful for concurrent or asynchronous interactions.

---

### Handling Errors

If a strategy fails or an error occurs, the server will respond with a message under the `error` topic:

```json
{
  "topic": "error",
  "payload": {
    "message": "some error message",
    "status": "status code",
    "trace": "if provided, a trace"
  },
  "ackId": "12345"
}
```

This message includes:

- A descriptive error `message`
- An optional `status` code
- A `trace` (for debugging, if available)
- The `ackId` to match the failed request

This structure allows the frontend to gracefully handle and display backend errors.
