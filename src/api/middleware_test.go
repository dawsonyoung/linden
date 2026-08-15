package api

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func Test_RequestID_HeaderAbsent_GeneratesAndEchoes(t *testing.T) {
	var captured string
	h := requestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = RequestIDFrom(r.Context())
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if captured == "" {
		t.Fatal("request ID not present in context")
	}
	if got := rec.Header().Get(requestIDHeader); got != captured {
		t.Errorf("echoed header = %q, want %q", got, captured)
	}
}

func Test_RequestID_ValidHeader_Preserved(t *testing.T) {
	const supplied = "client-supplied_123"

	var captured string
	h := requestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = RequestIDFrom(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set(requestIDHeader, supplied)
	h.ServeHTTP(httptest.NewRecorder(), req)

	if captured != supplied {
		t.Errorf("request ID = %q, want %q", captured, supplied)
	}
}

func Test_RequestID_UnsafeHeader_Replaced(t *testing.T) {
	tests := []struct {
		name     string
		supplied string
	}{
		{"newline injection", "abc\ninjected=true"},
		{"carriage return", "abc\r\nfake"},
		{"json breakout", `abc","forged":"yes`},
		{"space", "abc def"},
		{"over length limit", strings.Repeat("a", maxRequestIDLen+1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var captured string
			h := requestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				captured = RequestIDFrom(r.Context())
			}))

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			req.Header.Set(requestIDHeader, tt.supplied)
			h.ServeHTTP(httptest.NewRecorder(), req)

			if captured == tt.supplied {
				t.Errorf("unsafe request ID was accepted: %q", captured)
			}
			if captured == "" {
				t.Error("no replacement request ID generated")
			}
		})
	}
}

func Test_AccessLog_RecordsRequestMetadata(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	h := requestID(accessLog(logger, okHandler()))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/health", nil))

	out := buf.String()
	for _, want := range []string{`"method":"GET"`, `"path":"/health"`, `"status":200`, "duration_ms", "request_id"} {
		if !strings.Contains(out, want) {
			t.Errorf("log missing %s\ngot: %s", want, out)
		}
	}
}

func Test_AccessLog_QueryString_NotLogged(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	h := requestID(accessLog(logger, okHandler()))
	req := httptest.NewRequest(http.MethodGet, "/health?token=s3cret&q=user+question", nil)
	h.ServeHTTP(httptest.NewRecorder(), req)

	out := buf.String()
	for _, forbidden := range []string{"s3cret", "token", "user+question"} {
		if strings.Contains(out, forbidden) {
			t.Errorf("access log leaked query content %q\ngot: %s", forbidden, out)
		}
	}
}

func Test_AccessLog_RecordsHandlerStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	h := accessLog(logger, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/health", nil))

	if !strings.Contains(buf.String(), `"status":418`) {
		t.Errorf("log did not record handler status\ngot: %s", buf.String())
	}
}
