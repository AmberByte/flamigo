# Logging

Flamigo uses Go's standard `log/slog` package for framework-internal logs.

Framework records are emitted with structured attributes so you can filter or route them without disabling all debug logging:

- `library=flamigo`
- `component=<package area>`
- `code=<stable event code>`

For example, the event listener type-mismatch debug record includes:

- `component=events`
- `code=unsupported_event_type`

## Use Your App's Default slog Logger

If your application already uses `slog`, the simplest setup is to configure the default logger once:

```go
package main

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

Flamigo will pick that up automatically.

## Use a Dedicated Logger for Flamigo

If you want Flamigo logs to use a different handler than the rest of your app, set a framework-scoped logger:

```go
package main

import (
	flamigo "github.com/amberbyte/flamigo"
	"log/slog"
	"os"
)

func init() {
	flamigo.SetLogger(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}
```

## Filter Specific Framework Records

To suppress a noisy debug record while keeping other debug logs, wrap your handler with `flamigo.NewFilterHandler`:

```go
package main

import (
	"context"
	flamigo "github.com/amberbyte/flamigo"
	"log/slog"
	"os"
)

func init() {
	handler := flamigo.NewFilterHandler(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		func(_ context.Context, record slog.Record) bool {
			if record.Message != "unsupported event type" {
				return true
			}

			allow := true
			record.Attrs(func(attr slog.Attr) bool {
				if attr.Key == "code" && attr.Value.String() == "unsupported_event_type" {
					allow = false
					return false
				}
				return true
			})
			return allow
		},
	)

	flamigo.SetLogger(slog.New(handler))
}
```

## If Your App Still Uses Logrus

Your application can continue using `logrus`. The clean approach is to bridge Flamigo's `slog` output into a `logrus` logger with a custom `slog.Handler`.

```go
package main

import (
	"context"
	flamigo "github.com/amberbyte/flamigo"
	"log/slog"

	"github.com/sirupsen/logrus"
)

type logrusHandler struct {
	entry *logrus.Entry
	attrs []slog.Attr
}

func newLogrusHandler(logger *logrus.Logger) slog.Handler {
	return &logrusHandler{entry: logrus.NewEntry(logger)}
}

func (h *logrusHandler) Enabled(_ context.Context, level slog.Level) bool {
	switch {
	case level >= slog.LevelError:
		return h.entry.Logger.IsLevelEnabled(logrus.ErrorLevel)
	case level >= slog.LevelWarn:
		return h.entry.Logger.IsLevelEnabled(logrus.WarnLevel)
	case level >= slog.LevelInfo:
		return h.entry.Logger.IsLevelEnabled(logrus.InfoLevel)
	default:
		return h.entry.Logger.IsLevelEnabled(logrus.DebugLevel)
	}
}

func (h *logrusHandler) Handle(_ context.Context, record slog.Record) error {
	fields := logrus.Fields{}

	for _, attr := range h.attrs {
		fields[attr.Key] = attr.Value.Any()
	}
	record.Attrs(func(attr slog.Attr) bool {
		fields[attr.Key] = attr.Value.Any()
		return true
	})

	entry := h.entry.WithFields(fields)

	switch {
	case record.Level >= slog.LevelError:
		entry.Error(record.Message)
	case record.Level >= slog.LevelWarn:
		entry.Warn(record.Message)
	case record.Level >= slog.LevelInfo:
		entry.Info(record.Message)
	default:
		entry.Debug(record.Message)
	}

	return nil
}

func (h *logrusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	combined := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	combined = append(combined, h.attrs...)
	combined = append(combined, attrs...)
	return &logrusHandler{entry: h.entry, attrs: combined}
}

func (h *logrusHandler) WithGroup(string) slog.Handler {
	return h
}

func init() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	flamigo.SetLogger(slog.New(newLogrusHandler(logger)))
}
```
