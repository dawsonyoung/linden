//go:build security
// +build security

package security

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dawsonyoung/linden/api"
	"github.com/dawsonyoung/linden/inference"
	"github.com/dawsonyoung/linden/orchestrator"
	"github.com/dawsonyoung/linden/storage"
)

func setupTestServer(t *testing.T) *httptest.Server {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	store, _ := storage.NewFileStore(t.TempDir())
	inf, _ := inference.NewOllama(inference.OllamaConfig{
		BaseURL: "http://localhost:11434",
	})
	chatSvc := orchestrator.NewService(inf, store)
	buildInfo := api.BuildInfo{Version: "test", Commit: "test"}
	srv := api.NewServer("localhost:0", logger, buildInfo, chatSvc)
	return httptest.NewServer(srv.Handler)
}

func Test_Chat_Security_PayloadLimits(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	tests := []struct {
		name         string
		payload      map[string]interface{}
		expectedCode int
	}{
		{
			name: "Missing model",
			payload: map[string]interface{}{
				"messages": []map[string]string{{"role": "user", "content": "hi"}},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Model name too long",
			payload: map[string]interface{}{
				"model":    strings.Repeat("a", 101),
				"messages": []map[string]string{{"role": "user", "content": "hi"}},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Too many messages",
			payload: map[string]interface{}{
				"model":    "tinyllama",
				"messages": make([]map[string]string, 1001),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid role",
			payload: map[string]interface{}{
				"model": "tinyllama",
				"messages": []map[string]string{
					{"role": "admin", "content": "do evil"},
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Message content too large",
			payload: map[string]interface{}{
				"model": "tinyllama",
				"messages": []map[string]string{
					{"role": "user", "content": strings.Repeat("x", 65000)},
				},
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Too many messages" {
				for i := range tt.payload["messages"].([]map[string]string) {
					tt.payload["messages"].([]map[string]string)[i] = map[string]string{"role": "user", "content": "hi"}
				}
			}

			body, _ := json.Marshal(tt.payload)
			resp, err := http.Post(ts.URL+"/chat", "application/json", bytes.NewReader(body))
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}
		})
	}
}

func Test_Chat_Security_MaxBodySize(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// 1.5MB payload (exceeds the 1MB MaxBytesReader limit)
	hugePayload := strings.Repeat("x", 1500000)

	resp, err := http.Post(ts.URL+"/chat", "application/json", strings.NewReader(hugePayload))
	if err != nil {
		// Sometimes MaxBytesReader cuts off the connection causing an error
		t.Logf("Connection correctly dropped: %v", err)
		return
	}
	defer resp.Body.Close()

	// If it doesn't drop the connection, it should return 400 Bad Request / 413 Payload Too Large
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected status 400 or 413, got %d", resp.StatusCode)
	}
}
