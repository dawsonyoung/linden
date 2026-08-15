package inference

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dawsonyoung/linden/errs"
)

func newTestClient(t *testing.T, h http.Handler) *Ollama {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)

	client, err := NewOllama(OllamaConfig{BaseURL: server.URL, ResponseTimeout: time.Second})
	if err != nil {
		t.Fatalf("NewOllama() error = %v, want nil", err)
	}
	return client
}

func chatRequest() Request {
	return Request{
		Model:    "tinyllama",
		Messages: []Message{{Role: RoleUser, Content: "hello"}},
	}
}

func Test_NewOllama_InvalidConfig_ReturnsInvalidArgument(t *testing.T) {
	tests := map[string]string{
		"empty base URL":     "",
		"missing scheme":     "localhost:11434",
		"unsupported scheme": "ftp://localhost:11434",
	}

	for name, baseURL := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewOllama(OllamaConfig{BaseURL: baseURL})
			if err == nil {
				t.Fatal("NewOllama() error = nil, want an error")
			}
			if got := errs.CodeOf(err); got != errs.InvalidArgument {
				t.Errorf("code = %q, want %q", got, errs.InvalidArgument)
			}
		})
	}
}

func Test_NewOllama_TrailingSlash_IsNormalized(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"models":[]}`))
	}))
	defer server.Close()

	client, err := NewOllama(OllamaConfig{BaseURL: server.URL + "/"})
	if err != nil {
		t.Fatalf("NewOllama() error = %v, want nil", err)
	}
	if _, err := client.ListModels(context.Background()); err != nil {
		t.Fatalf("ListModels() error = %v, want nil", err)
	}
	if gotPath != "/api/tags" {
		t.Errorf("request path = %q, want %q", gotPath, "/api/tags")
	}
}

func Test_ChatStream_SendsStreamingRequest(t *testing.T) {
	var body struct {
		Model    string `json:"model"`
		Stream   bool   `json:"stream"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := decodeJSON(r, &body); err != nil {
			t.Errorf("decode request: %v", err)
		}
		_, _ = w.Write([]byte(`{"done":true}` + "\n"))
	}))

	if _, err := client.ChatStream(context.Background(), chatRequest(), func(Chunk) error { return nil }); err != nil {
		t.Fatalf("ChatStream() error = %v, want nil", err)
	}

	if !body.Stream {
		t.Error("request did not set stream=true")
	}
	if body.Model != "tinyllama" {
		t.Errorf("model = %q, want %q", body.Model, "tinyllama")
	}
	if len(body.Messages) != 1 || body.Messages[0].Role != "user" {
		t.Errorf("messages = %+v, want one user message", body.Messages)
	}
}

func Test_ChatStream_MalformedFrame_ReturnsInternal(t *testing.T) {
	tests := map[string]string{
		"not json":         "this is not json\n",
		"truncated object": `{"message":{"content":"hi"` + "\n",
		"json array":       `["unexpected"]` + "\n",
	}

	for name, payload := range tests {
		t.Run(name, func(t *testing.T) {
			client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(payload))
			}))

			_, err := client.ChatStream(context.Background(), chatRequest(), func(Chunk) error { return nil })
			if err == nil {
				t.Fatal("ChatStream() error = nil, want an error")
			}
			if got := errs.CodeOf(err); got != errs.Internal {
				t.Errorf("code = %q, want %q", got, errs.Internal)
			}
		})
	}
}

func Test_ChatStream_StreamEndsWithoutTerminalFrame_ReturnsInternal(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"message":{"content":"partial"}}` + "\n"))
	}))

	_, err := client.ChatStream(context.Background(), chatRequest(), func(Chunk) error { return nil })
	if err == nil {
		t.Fatal("ChatStream() error = nil, want an error")
	}
	if got := errs.CodeOf(err); got != errs.Internal {
		t.Errorf("code = %q, want %q", got, errs.Internal)
	}
}

func Test_ChatStream_ProviderErrorFrame_ReturnsInternalWithoutEchoingIt(t *testing.T) {
	const secret = "user asked about their medical history"

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"error":"failed generating for: ` + secret + `"}` + "\n"))
	}))

	_, err := client.ChatStream(context.Background(), chatRequest(), func(Chunk) error { return nil })
	if err == nil {
		t.Fatal("ChatStream() error = nil, want an error")
	}
	// The provider's text may echo the request, so it must not reach the message.
	if strings.Contains(err.Error(), secret) {
		t.Errorf("error message echoed provider content: %q", err.Error())
	}
}

func Test_ChatStream_SkipsEmptyAndBlankFrames(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("\n" +
			`{"message":{"content":""}}` + "\n" +
			"   \n" +
			`{"message":{"content":"real"}}` + "\n" +
			`{"done":true}` + "\n"))
	}))

	var got []string
	if _, err := client.ChatStream(context.Background(), chatRequest(), func(c Chunk) error {
		got = append(got, c.Text)
		return nil
	}); err != nil {
		t.Fatalf("ChatStream() error = %v, want nil", err)
	}

	if len(got) != 1 || got[0] != "real" {
		t.Errorf("chunks = %q, want exactly [real]", got)
	}
}

func Test_StatusMapping(t *testing.T) {
	tests := []struct {
		status int
		want   errs.Code
	}{
		{http.StatusNotFound, errs.NotFound},
		{http.StatusBadRequest, errs.InvalidArgument},
		{http.StatusTooManyRequests, errs.ResourceExhausted},
		{http.StatusInternalServerError, errs.Unavailable},
		{http.StatusBadGateway, errs.Unavailable},
		{http.StatusTeapot, errs.Internal},
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
			}))

			_, err := client.ChatStream(context.Background(), chatRequest(), func(Chunk) error { return nil })
			if err == nil {
				t.Fatal("ChatStream() error = nil, want an error")
			}
			if got := errs.CodeOf(err); got != tt.want {
				t.Errorf("code = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_ListModels_MalformedPayload_ReturnsInternal(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"models":`))
	}))

	_, err := client.ListModels(context.Background())
	if err == nil {
		t.Fatal("ListModels() error = nil, want an error")
	}
	if got := errs.CodeOf(err); got != errs.Internal {
		t.Errorf("code = %q, want %q", got, errs.Internal)
	}
}

func Test_ListModels_EmptyList_ReturnsEmptySlice(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"models":[]}`))
	}))

	got, err := client.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels() error = %v, want nil", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d models, want 0", len(got))
	}
}

func Test_AlreadyCanceledContext_SkipsProviderCall(t *testing.T) {
	called := false
	client := newTestClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := client.ListModels(ctx); !errors.Is(err, context.Canceled) {
		t.Errorf("ListModels() error = %v, want context.Canceled", err)
	}

	res, err := client.ChatStream(ctx, chatRequest(), func(Chunk) error { return nil })
	if !errors.Is(err, context.Canceled) {
		t.Errorf("ChatStream() error = %v, want context.Canceled", err)
	}
	if res.FinishReason != FinishCanceled {
		t.Errorf("FinishReason = %q, want %q", res.FinishReason, FinishCanceled)
	}
	if called {
		t.Error("provider was contacted despite a canceled context")
	}
}

func Test_FrameExceedingLimit_ReturnsError(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		huge := strings.Repeat("a", maxFrameBytes+1024)
		_, _ = w.Write([]byte(`{"message":{"content":"` + huge + `"}}` + "\n"))
	}))

	_, err := client.ChatStream(context.Background(), chatRequest(), func(Chunk) error { return nil })
	if err == nil {
		t.Fatal("ChatStream() error = nil, want an error for an oversized frame")
	}
}

func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
