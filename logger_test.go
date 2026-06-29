package flamigo

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

func TestNewFilterHandler(t *testing.T) {
	var output bytes.Buffer
	handler := NewFilterHandler(
		slog.NewTextHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}),
		func(_ context.Context, record slog.Record) bool {
			return record.Message != "drop"
		},
	)

	logger := slog.New(handler)
	logger.Debug("keep", slog.String("code", "keep"))
	logger.Debug("drop", slog.String("code", "drop"))

	logs := output.String()
	if logs == "" {
		t.Fatal("expected at least one log line")
	}
	if bytes.Contains(output.Bytes(), []byte("drop")) {
		t.Fatal("expected matching record to be filtered")
	}
}

func TestSetLoggerNilRestoresDefault(t *testing.T) {
	previous := Logger()
	defer SetLogger(previous)

	SetLogger(nil)
	if Logger() == nil {
		t.Fatal("expected logger to remain non-nil")
	}
}
