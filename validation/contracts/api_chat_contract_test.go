//go:build contracts

package contracts

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/dawsonyoung/linden/api"
	"github.com/dawsonyoung/linden/orchestrator"
)

func doJSON(t *testing.T, h http.Handler, method, path string, body any, headers map[string]string) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	reqHeaders := map[string]string{"Content-Type": "application/json"}
	for k, v := range headers {
		reqHeaders[k] = v
	}
	return do(t, h, method, path, bytes.NewReader(b), reqHeaders)
}

func Test_Chat_Post_EmptyMessages_Returns400(t *testing.T) {
	h := newHandler(t, api.BuildInfo{}, &scriptedChatService{})

	reqBody := map[string]any{
		"model":    "llama2",
		"messages": []any{},
	}
	resp := doJSON(t, h, http.MethodPost, "/chat", reqBody, nil)
	requireStatus(t, resp, http.StatusBadRequest)
	body := requireJSONObject(t, resp)

	if got := requireStringField(t, body, "error"); got == "" {
		t.Error("error field is missing or empty")
	}
}

func Test_Chat_Post_UnknownModel_Returns404(t *testing.T) {
	s := orchestratorScript{
		Models: []orchestrator.Model{{Name: "llama2"}},
	}
	h := newHandler(t, api.BuildInfo{}, &scriptedChatService{script: s})

	reqBody := map[string]any{
		"model":    "gpt-4",
		"messages": []map[string]string{{"role": "user", "content": "hi"}},
	}
	resp := doJSON(t, h, http.MethodPost, "/chat", reqBody, nil)
	requireStatus(t, resp, http.StatusNotFound)
}

func Test_Chat_Post_SSE_ReturnsStream(t *testing.T) {
	s := orchestratorScript{
		Models: []orchestrator.Model{{Name: "llama2"}},
		Chunks: []string{"Hello", " ", "World"},
	}
	h := newHandler(t, api.BuildInfo{}, &scriptedChatService{script: s})

	reqBody := map[string]any{
		"model":    "llama2",
		"messages": []map[string]string{{"role": "user", "content": "hi"}},
		"stream":   true,
	}
	resp := doJSON(t, h, http.MethodPost, "/chat", reqBody, nil)
	requireStatus(t, resp, http.StatusOK)

	if got := requireHeader(t, resp, "Content-Type"); !strings.HasPrefix(got, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want text/event-stream", got)
	}

	scanner := bufio.NewScanner(resp.Body)
	var events []string
	var currentEvent []string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if len(currentEvent) > 0 {
				events = append(events, strings.Join(currentEvent, "\n"))
				currentEvent = nil
			}
			continue
		}
		currentEvent = append(currentEvent, line)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan SSE stream: %v", err)
	}

	if len(events) < 2 {
		t.Fatalf("expected multiple SSE events, got %d", len(events))
	}

	lastEvent := events[len(events)-1]
	if !strings.Contains(lastEvent, "event: done") && !strings.Contains(lastEvent, "event: error") {
		t.Errorf("terminal event missing, last event was:\n%s", lastEvent)
	}
}
