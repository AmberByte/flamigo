# WebSocket (Realtime)

Flamigo's realtime WebSocket adapter is built on **WebSockets**, allowing you to send or push updates from your backend to connected frontend clients instantly and efficiently.

It integrates directly with the **Realtime Event Bus**, enabling seamless reaction to domain events and bidirectional communication with frontend applications.

::: warning
WebSocket support can only be enabled when initializing a new project because it requires specific scaffolding to be generated.
:::

---

## Setting Up WebSocket Communication

Follow these steps to enable and use WebSocket communication in your Flamigo project:

1. **Enable WebSocket Support**:
   Ensure the WebSocket module is enabled during project initialization. If not, you may need to reinitialize the project.

2. **Define WebSocket Strategies**:
   Create strategies to handle WebSocket messages. For example:
   ```go
   func MyWebSocketStrategy(ctx flamigo.Context, payload map[string]interface{}) error {
       // Handle the WebSocket message
       return nil
   }
   ```

3. **Subscribe Clients to Topics**:
   In the `client.go` file, subscribe WebSocket clients to topics:
   ```go
   client.Subscribe("app::strategy", MyWebSocketStrategy)
   ```

4. **Frontend Integration**:
   Use a WebSocket client library (e.g., `native WebSocket` or `socket.io`) to connect to the backend:
   ```javascript
   const socket = new WebSocket("ws://your-backend-url");
   socket.onmessage = (event) => {
       console.log("Message received:", event.data);
   };
   ```

5. **Test the Connection**:
   Send a test message from the frontend and verify the backend processes it correctly.

---

## WebSocket Actor

In line with Flamigo’s actor-based architecture, each WebSocket connection is treated as an **actor**. This allows strategies and domain logic to understand and respond to *who* initiated a given action via WebSocket.

The WebSocket adapter provides its own actor type to represent the client.

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

If the authentication plugin is enabled, the WebSocket adapter registers a dedicated **authentication strategy**:

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

The generated WebSocket scaffold under `internal/adapters/websocket` keeps auth, connection lifecycle, and event projection in your app code, but delegates reusable transport mechanics like command decoding, strategy dispatch, and response encoding to `transport/websocket`.

The WebSocket adapter emits lifecycle events, including:

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

To receive a response to a specific request, include an `ack` key in the message:

```json
{
  "topic": "app::strategy",
  "payload": {},
  "ack": "12345"
}
```

The response will be returned with the same `ack` value:

```json
{
  "topic": "app::strategy",
  "payload": {
    "foo": "bar"
  },
  "ack": "12345"
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
    "type": "server_error",
    "trace": "if provided, a trace"
  },
  "ack": "12345"
}
```

This message includes:

- A descriptive error `message`
- An optional generic error `type`, such as `bad_request`, `not_found`, or `server_error`
- A `trace` (for debugging, if available)
- The `ack` value to match the failed request

This structure allows the frontend to gracefully handle and display backend errors.
