//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dawsonyoung/linden/api"
	"github.com/dawsonyoung/linden/inference"
	"github.com/dawsonyoung/linden/orchestrator"
	"github.com/dawsonyoung/linden/storage"
)

func newSessionIntegrationHarness(t *testing.T, ollamaHandler http.Handler) *httptest.Server {
	t.Helper()

	providerSrv := httptest.NewServer(ollamaHandler)
	t.Cleanup(providerSrv.Close)

	inf, err := inference.NewOllama(inference.OllamaConfig{
		BaseURL:         providerSrv.URL,
		ResponseTimeout: 2 * time.Second,
	})
	if err != nil {
		t.Fatalf("build inference: %v", err)
	}

	tmpDir := t.TempDir()
	store, err := storage.NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("build storage: %v", err)
	}

	chatService := orchestrator.NewService(inf, store)

	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	apiSrv := api.NewServer("localhost:0", logger, api.BuildInfo{}, chatService)
	testSrv := httptest.NewServer(apiSrv.Handler)
	t.Cleanup(testSrv.Close)

	return testSrv
}

func Test_Integration_Session_MultiTurnMemory(t *testing.T) {
	var capturedRequests []inference.Request

	ollama := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req inference.Request
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &req)
		capturedRequests = append(capturedRequests, req)

		w.Header().Set("Content-Type", "application/x-ndjson")
		chunk1 := `{"model":"tinyllama","message":{"role":"assistant","content":"hello"},"done":true}` + "\n"
		w.Write([]byte(chunk1))
	})

	apiSrv := newSessionIntegrationHarness(t, ollama)

	// Turn 1
	req1Body := map[string]any{
		"model":     "tinyllama",
		"sessionId": "session-123",
		"messages": []map[string]string{
			{"role": "user", "content": "my first message"},
		},
	}
	b1, _ := json.Marshal(req1Body)
	req1, _ := http.NewRequest(http.MethodPost, apiSrv.URL+"/chat", bytes.NewReader(b1))
	req1.Header.Set("Content-Type", "application/json")

	resp1, err := apiSrv.Client().Do(req1)
	if err != nil {
		t.Fatalf("do chat request 1: %v", err)
	}
	io.Copy(io.Discard, resp1.Body)
	resp1.Body.Close()

	if len(capturedRequests) != 1 {
		t.Fatalf("expected 1 inference request, got %d", len(capturedRequests))
	}
	if len(capturedRequests[0].Messages) != 1 {
		t.Errorf("turn 1 inference request should have 1 message, got %d", len(capturedRequests[0].Messages))
	}

	// Turn 2
	req2Body := map[string]any{
		"model":     "tinyllama",
		"sessionId": "session-123",
		"messages": []map[string]string{
			{"role": "user", "content": "my second message"},
		},
	}
	b2, _ := json.Marshal(req2Body)
	req2, _ := http.NewRequest(http.MethodPost, apiSrv.URL+"/chat", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := apiSrv.Client().Do(req2)
	if err != nil {
		t.Fatalf("do chat request 2: %v", err)
	}
	io.Copy(io.Discard, resp2.Body)
	resp2.Body.Close()

	if len(capturedRequests) != 2 {
		t.Fatalf("expected 2 inference requests, got %d", len(capturedRequests))
	}

	messages := capturedRequests[1].Messages
	if len(messages) != 3 {
		t.Errorf("turn 2 inference request should have 3 messages, got %d", len(messages))
	} else {
		if messages[0].Content != "my first message" {
			t.Errorf("expected msg 0 to be 'my first message', got %q", messages[0].Content)
		}
		if messages[1].Content != "hello" {
			t.Errorf("expected msg 1 to be 'hello', got %q", messages[1].Content)
		}
		if messages[2].Content != "my second message" {
			t.Errorf("expected msg 2 to be 'my second message', got %q", messages[2].Content)
		}
	}
}
