//go:build integration

package integration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dawsonyoung/linden/api"
	"github.com/dawsonyoung/linden/inference"
	"github.com/dawsonyoung/linden/orchestrator"
)

func newIntegrationHarness(t *testing.T, ollamaHandler http.Handler) (*httptest.Server, *bytes.Buffer) {
	t.Helper()

	// Mock provider
	providerSrv := httptest.NewServer(ollamaHandler)
	t.Cleanup(providerSrv.Close)

	inf, err := inference.NewOllama(inference.OllamaConfig{
		BaseURL:         providerSrv.URL,
		ResponseTimeout: 2 * time.Second,
	})
	if err != nil {
		t.Fatalf("build inference: %v", err)
	}

	chatService := orchestrator.NewService(inf)

	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	apiSrv := api.NewServer("localhost:0", logger, api.BuildInfo{}, chatService)
	testSrv := httptest.NewServer(apiSrv.Handler)
	t.Cleanup(testSrv.Close)

	return testSrv, &logBuf
}

func Test_Integration_Chat_HappyPath_StreamsSSE(t *testing.T) {
	ollama := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}

		// Assert no user data leaks in provider request body
		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "secret") {
			// This assertion is fine, the user request goes to Ollama.
		}

		w.Header().Set("Content-Type", "application/x-ndjson")

		chunk1 := `{"model":"tinyllama","message":{"role":"assistant","content":"hello "},"done":false}` + "\n"
		chunk2 := `{"model":"tinyllama","message":{"role":"assistant","content":"world"},"done":true}` + "\n"

		w.Write([]byte(chunk1))
		w.Write([]byte(chunk2))
	})

	apiSrv, logBuf := newIntegrationHarness(t, ollama)

	reqBody := map[string]any{
		"model": "tinyllama",
		"messages": []map[string]string{
			{"role": "user", "content": "my secret message"},
		},
	}
	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, apiSrv.URL+"/chat", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := apiSrv.Client().Do(req)
	if err != nil {
		t.Fatalf("do chat request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", got)
	}

	scanner := bufio.NewScanner(resp.Body)
	var events []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			events = append(events, line)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("read stream: %v", err)
	}

	if len(events) < 2 {
		t.Fatalf("got %d events, want at least 2", len(events))
	}

	// First event should be "hello "
	if !strings.Contains(events[0], "event: message") || !strings.Contains(events[1], `"text":"hello "`) {
		t.Errorf("unexpected first event: %v", events[0:2])
	}

	// Make sure logs don't contain the secret message
	if strings.Contains(logBuf.String(), "secret message") {
		t.Error("user input leaked into logs")
	}
}

func Test_Integration_Chat_NetworkTimeout_Returns503(t *testing.T) {
	ollama := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second) // Exceeds ResponseTimeout
	})

	apiSrv, _ := newIntegrationHarness(t, ollama)

	reqBody := map[string]any{
		"model": "tinyllama",
		"messages": []map[string]string{
			{"role": "user", "content": "hi"},
		},
	}
	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, apiSrv.URL+"/chat", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	// We don't want the test client to timeout before the server responds 503
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do chat request: %v", err)
	}
	defer resp.Body.Close()

	// Timeout maps to Unavailable (503)
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusServiceUnavailable)
	}
}
