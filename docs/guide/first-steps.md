# First Steps

## 1. Install the CLI

```bash
go install github.com/amberbyte/flamigo/tools/flamigo@latest
```

## 2. Scaffold a project

```bash
flamigo init
```

The wizard asks for:

- your Go module path
- optional features to enable

Current optional features include:

- `auth`
- `transport_http`
- `transport_websocket`

## 3. Understand the generated structure

The scaffold gives you a project shaped around:

- `internal/domains` for business logic
- `internal/api` for strategy registration
- `internal/adapters` for transport and external integrations
- `pkg` for app-owned reusable helpers

## 4. Wire your application at startup

The generated `cmd/main.go` wires packages through dependency injection during startup.

It also configures application logging with `log/slog`, which Flamigo uses internally for framework logs as well.

In practice, the typical flow is:

1. initialize core framework pieces
2. initialize domain infrastructure
3. register API strategies
4. initialize adapters
5. start transport servers

If you need to route or filter framework logs separately, see the [Logging](/guide/logging) guide.

## 5. Add your first strategy

The next useful step is usually to add a strategy in `internal/api`, register it in the app registry, and then expose it through an adapter such as HTTP or WebSocket.

From there, grow domain logic inside `internal/domains`, not inside the transport layer.
