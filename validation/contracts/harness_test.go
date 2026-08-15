//go:build contracts

// Shared assertions for contract tests. These live in the package that uses
// them rather than a separate harness package; a second consumer would justify
// extracting them, one does not.
package contracts

import (
	"encoding/json"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dawsonyoung/linden/api"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}

// newHandler returns the server's HTTP handler as a client would reach it, with
// middleware applied. Contract tests exercise the published surface, so they
// must not invoke handlers directly.
func newHandler(t *testing.T, build api.BuildInfo) http.Handler {
	t.Helper()
	srv := api.NewServer(":0", discardLogger(), build)
	if srv.Handler == nil {
		t.Fatal("NewServer returned a server with no handler")
	}
	return srv.Handler
}

func do(t *testing.T, h http.Handler, method, path string, headers map[string]string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Result()
}

func requireStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		t.Fatalf("status = %d, want %d", resp.StatusCode, want)
	}
}

func requireHeader(t *testing.T, resp *http.Response, key string) string {
	t.Helper()
	v := resp.Header.Get(key)
	if v == "" {
		t.Fatalf("response is missing header %s", key)
	}
	return v
}

// requireJSONObject decodes a body the published contract describes as a JSON
// object. The media type is compared without parameters: a charset would
// satisfy every client and must not fail the contract.
func requireJSONObject(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()

	mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("parse Content-Type %q: %v", resp.Header.Get("Content-Type"), err)
	}
	if mediaType != "application/json" {
		t.Fatalf("media type = %q, want %q", mediaType, "application/json")
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body as JSON object: %v", err)
	}
	return body
}

func requireStringField(t *testing.T, body map[string]any, field string) string {
	t.Helper()
	raw, ok := body[field]
	if !ok {
		t.Fatalf("response has no %q field; got %v", field, fieldNames(body))
	}
	s, ok := raw.(string)
	if !ok {
		t.Fatalf("field %q = %T, want string", field, raw)
	}
	if s == "" {
		t.Fatalf("field %q is empty", field)
	}
	return s
}

func fieldNames(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
