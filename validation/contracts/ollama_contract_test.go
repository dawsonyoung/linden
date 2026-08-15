//go:build contracts

package contracts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dawsonyoung/linden/inference"
)

// Runs the same conformance suite A.2 defined, against the real adapter backed
// by a stub provider. If the suite encoded assumptions only an in-memory double
// could satisfy, this is where that shows.
func Test_OllamaClient_SatisfiesClientContract(t *testing.T) {
	runClientContract(t, newOllamaUnderScript)
}

func newOllamaUnderScript(t *testing.T, s script) inference.Client {
	t.Helper()

	baseURL := stubProviderURL(t, s)

	client, err := inference.NewOllama(inference.OllamaConfig{
		BaseURL: baseURL,
		// Short so the stalled-provider scenario fails fast rather than
		// consuming the suite's context budget.
		ResponseTimeout: 150 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewOllama() error = %v, want nil", err)
	}
	return client
}

// stubProviderURL serves enough of the Ollama API to exercise the contract.
func stubProviderURL(t *testing.T, s script) string {
	t.Helper()

	if s.Fail == failUnreachable {
		// A listener that is closed immediately: connecting is refused.
		dead := httptest.NewServer(http.NotFoundHandler())
		url := dead.URL
		dead.Close()
		return url
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/tags", func(w http.ResponseWriter, r *http.Request) {
		if s.Fail == failTimeout {
			stall(r)
			return
		}
		type wireModel struct {
			Name string `json:"name"`
		}
		payload := struct {
			Models []wireModel `json:"models"`
		}{}
		for _, m := range s.Models {
			payload.Models = append(payload.Models, wireModel{Name: m.Name})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	})

	mux.HandleFunc("POST /api/chat", func(w http.ResponseWriter, r *http.Request) {
		if s.Fail == failTimeout {
			stall(r)
			return
		}

		var req struct {
			Model string `json:"model"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !scriptServes(s, req.Model) {
			// Ollama reports an unknown model as 404.
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"model not found"}`))
			return
		}

		w.Header().Set("Content-Type", "application/x-ndjson")
		w.WriteHeader(http.StatusOK)
		flusher, _ := w.(http.Flusher)

		for _, text := range s.Chunks {
			frame := fmt.Sprintf(`{"model":%q,"done":false,"message":{"role":"assistant","content":%q}}`, req.Model, text)
			if _, err := fmt.Fprintln(w, frame); err != nil {
				return
			}
			if flusher != nil {
				flusher.Flush()
			}
		}
		_, _ = fmt.Fprintf(w, `{"model":%q,"done":true,"done_reason":"stop"}`+"\n", req.Model)
		if flusher != nil {
			flusher.Flush()
		}
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server.URL
}

// stall blocks without writing response headers, which is what a provider that
// has accepted a connection but cannot answer looks like.
//
// Bounded deliberately: the client aborts on ResponseHeaderTimeout without
// necessarily closing the connection, so waiting only on the request context
// would leave this handler running and block server shutdown.
func stall(r *http.Request) {
	select {
	case <-r.Context().Done():
	case <-time.After(time.Second):
	}
}

func scriptServes(s script, model string) bool {
	for _, m := range s.Models {
		if m.Name == model {
			return true
		}
	}
	return false
}
