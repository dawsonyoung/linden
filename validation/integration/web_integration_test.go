//go:build integration
// +build integration

package integration

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/dawsonyoung/linden/api"
	"github.com/dawsonyoung/linden/inference"
	"github.com/dawsonyoung/linden/orchestrator"
	"github.com/dawsonyoung/linden/storage"
)

func Test_WebUI_Integration_ServesStaticAssets(t *testing.T) {
	// Create test dependencies
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Use an in-memory store and a dummy inference client since we just test UI serving
	store, _ := storage.NewFileStore(t.TempDir())
	inf, _ := inference.NewOllama(inference.OllamaConfig{
		BaseURL: "http://localhost:11434",
	})
	chatSvc := orchestrator.NewService(inf, store)

	// Create the API server
	buildInfo := api.BuildInfo{Version: "test", Commit: "test"}
	srv := api.NewServer("localhost:0", logger, buildInfo, chatSvc)

	// Extract the actual handler from the server
	handler := srv.Handler

	// Spin up an httptest server using our handler
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Perform a GET request to the root URL
	client := ts.Client()
	client.Timeout = 2 * time.Second

	resp, err := client.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("Failed to GET /: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type text/html; charset=utf-8, got %q", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	bodyStr := string(body)
	if !contains(bodyStr, "Linden") && !contains(bodyStr, "sveltekit") {
		t.Errorf("Response body does not appear to contain SvelteKit static UI markers. Body: %s", bodyStr)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && len(s) >= len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && contains(s[1:], substr))
}

func Test_WebUI_Integration_BrowserE2E(t *testing.T) {
	// Require node in path
	if _, err := os.Stat("../../src/web/test-browser.mjs"); os.IsNotExist(err) {
		t.Skip("test-browser.mjs not found, skipping browser E2E test")
	}

	targetURL := os.Getenv("LINDEN_TEST_URL")

	if targetURL == "" {
		logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
		store, _ := storage.NewFileStore(t.TempDir())
		inf, _ := inference.NewOllama(inference.OllamaConfig{
			BaseURL: "http://localhost:11434",
		})
		chatSvc := orchestrator.NewService(inf, store)

		buildInfo := api.BuildInfo{Version: "test", Commit: "test"}
		srv := api.NewServer("localhost:0", logger, buildInfo, chatSvc)

		ts := httptest.NewServer(srv.Handler)
		defer ts.Close()
		targetURL = ts.URL
	} else {
		t.Logf("Running browser E2E test against external URL: %s", targetURL)
	}

	cmd := exec.Command("node", "../../src/web/test-browser.mjs", targetURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Browser E2E test failed. Ensure no JS errors occurred. Error: %v", err)
	}
}

func Test_WebUI_Integration_BrowserE2E_ErrorHandling(t *testing.T) {
	// Require node in path
	if _, err := os.Stat("../../src/web/test-browser-error.mjs"); os.IsNotExist(err) {
		t.Skip("test-browser-error.mjs not found, skipping error E2E test")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	store, _ := storage.NewFileStore(t.TempDir())
	
	// Use an Ollama backend that does NOT exist to force a connection error
	inf, _ := inference.NewOllama(inference.OllamaConfig{
		BaseURL: "http://localhost:59999", // Invalid port
	})
	chatSvc := orchestrator.NewService(inf, store)

	buildInfo := api.BuildInfo{Version: "test", Commit: "test"}
	srv := api.NewServer("localhost:0", logger, buildInfo, chatSvc)

	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	cmd := exec.Command("node", "../../src/web/test-browser-error.mjs", ts.URL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Browser error E2E test failed. UI did not handle the error gracefully. Error: %v", err)
	}
}
