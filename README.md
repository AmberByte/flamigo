# Flamigo

Hexagonal, event-driven Go backend framework.

> [!IMPORTANT]
> Flamigo is still pre-1.0. The architecture is already usable, but APIs, templates, and tooling may still change between releases.

![flamigo](docs/public/logo.png)

## What Flamigo Is

Flamigo is a backend framework for Go projects that want:

- domain-first structure
- hexagonal architecture
- event-driven coordination between domains
- transport adapters for HTTP and WebSocket-style applications
- lightweight dependency injection for startup wiring

It is designed for applications like game backends, internal platforms, modular monoliths, and services where domain boundaries matter more than raw CRUD scaffolding.

## What Flamigo Is Not

Flamigo is not trying to be:

- a full ORM or database framework
- a batteries-included auth product
- a frontend framework
- a giant “everything included” web framework

The framework gives you structure and reusable primitives. Application policy, persistence details, and adapter behavior still live in your app.

## Core Ideas

- `domains` hold business logic and domain contracts
- `events` carry domain events across the application
- `strategies` provide application actions that transports can invoke
- `adapters` connect the system to the outside world
- `injection` wires everything together at startup

## Getting Started

Install the CLI:

```bash
go install github.com/amberbyte/flamigo/tools/flamigo@latest
```

Create a project:

```bash
flamigo init
```

The wizard scaffolds a project and lets you enable optional features like:

- `auth`
- `transport_http`
- `transport_websocket`

## Documentation

Documentation and guides are available at:

[Documentation](https://flamigo.amberbyte.dev)
[GoDoc](https://pkg.go.dev/github.com/amberbyte/flamigo)

## Logging

Flamigo framework internals use Go's standard `log/slog` package and scope framework records with `library=flamigo`.

If you already use `slog` in your app, setting the default logger is enough:

```go
import (
	"log/slog"
	"os"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}
```

If you want Flamigo to use a dedicated logger instead of the global default, set it explicitly:

```go
import (
	flamigo "github.com/amberbyte/flamigo"
	"log/slog"
	"os"
)

func init() {
	flamigo.SetLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}
```

To suppress one noisy framework debug record while keeping other debug logs, wrap the handler with `flamigo.NewFilterHandler` and filter by structured attributes such as `component` and `code`.

If your app still uses `logrus`, bridge Flamigo's `slog` logger into a custom `slog.Handler` and pass it to `flamigo.SetLogger(...)`. See [docs/guide/logging.md](docs/guide/logging.md) for a full example.

## Contributing

Issues, ideas, and pull requests are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT
