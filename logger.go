package flamigo

import (
	"context"
	"log/slog"
	"sync"
)

var (
	loggerMu sync.RWMutex
	logger   = defaultLogger()
)

func defaultLogger() *slog.Logger {
	return slog.Default().With("library", "flamigo")
}

func Logger() *slog.Logger {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	return logger
}

func SetLogger(l *slog.Logger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	if l == nil {
		logger = defaultLogger()
		return
	}
	logger = l.With("library", "flamigo")
}

type Filter func(ctx context.Context, record slog.Record) bool

func NewFilterHandler(next slog.Handler, filter Filter) slog.Handler {
	if next == nil {
		next = slog.DiscardHandler
	}
	if filter == nil {
		filter = func(context.Context, slog.Record) bool { return true }
	}
	return &filterHandler{
		next:   next,
		filter: filter,
	}
}

type filterHandler struct {
	next   slog.Handler
	filter Filter
}

func (h *filterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *filterHandler) Handle(ctx context.Context, record slog.Record) error {
	if !h.filter(ctx, record) {
		return nil
	}
	return h.next.Handle(ctx, record)
}

func (h *filterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &filterHandler{
		next:   h.next.WithAttrs(attrs),
		filter: h.filter,
	}
}

func (h *filterHandler) WithGroup(name string) slog.Handler {
	return &filterHandler{
		next:   h.next.WithGroup(name),
		filter: h.filter,
	}
}
