package project

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

type templateTestData struct {
	ProjectModulePath string
	ProjectName       string
	Features          []string
}

func (t templateTestData) HasFeature(feature string) bool {
	return slices.Contains(t.Features, feature)
}

func TestInitializeProjectFiles_HTTPFeatureGate(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		root := t.TempDir()

		err := InitializeProjectFiles(root, templateTestData{
			ProjectModulePath: "github.com/example/enabled",
			ProjectName:       "enabled",
			Features:          []string{"transport_http"},
		})
		if err != nil {
			t.Fatalf("InitializeProjectFiles returned error: %v", err)
		}

		assertFileExists(t, filepath.Join(root, "internal", "adapters", "http", "server.go"))
		assertFileExists(t, filepath.Join(root, "internal", "adapters", "http", "actor.go"))
	})

	t.Run("disabled", func(t *testing.T) {
		root := t.TempDir()

		err := InitializeProjectFiles(root, templateTestData{
			ProjectModulePath: "github.com/example/disabled",
			ProjectName:       "disabled",
		})
		if err != nil {
			t.Fatalf("InitializeProjectFiles returned error: %v", err)
		}

		assertFileMissing(t, filepath.Join(root, "internal", "adapters", "http", "server.go"))
		assertFileMissing(t, filepath.Join(root, "internal", "adapters", "http", "actor.go"))
	})
}

func TestInitializeProjectFiles_WebsocketConditionFiles(t *testing.T) {
	t.Run("without auth", func(t *testing.T) {
		root := t.TempDir()

		err := InitializeProjectFiles(root, templateTestData{
			ProjectModulePath: "github.com/example/ws-no-auth",
			ProjectName:       "ws-no-auth",
			Features:          []string{"transport_websocket"},
		})
		if err != nil {
			t.Fatalf("InitializeProjectFiles returned error: %v", err)
		}

		assertFileMissing(t, filepath.Join(root, "internal", "adapters", "websocket", "auth.go"))

		server := readFile(t, filepath.Join(root, "internal", "adapters", "websocket", "server.go"))
		if strings.Contains(server, "httptransport") {
			t.Fatalf("expected websocket server without HTTP transport import, got:\n%s", server)
		}

		actor := readFile(t, filepath.Join(root, "internal", "adapters", "websocket", "actor.go"))
		if strings.Contains(actor, "auth.UserActor") {
			t.Fatalf("expected websocket actor without auth references, got:\n%s", actor)
		}
	})

	t.Run("with auth", func(t *testing.T) {
		root := t.TempDir()

		err := InitializeProjectFiles(root, templateTestData{
			ProjectModulePath: "github.com/example/ws-auth",
			ProjectName:       "ws-auth",
			Features:          []string{"transport_websocket", "auth"},
		})
		if err != nil {
			t.Fatalf("InitializeProjectFiles returned error: %v", err)
		}

		assertFileExists(t, filepath.Join(root, "internal", "adapters", "websocket", "auth.go"))

		actor := readFile(t, filepath.Join(root, "internal", "adapters", "websocket", "actor.go"))
		if !strings.Contains(actor, "auth.UserActor") {
			t.Fatalf("expected websocket actor auth references, got:\n%s", actor)
		}
	})
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %s: %v", path, err)
	}
}

func assertFileMissing(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected file to be absent: %s", path)
	} else if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error for %s: %v", path, err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return string(content)
}
