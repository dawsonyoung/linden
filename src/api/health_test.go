package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HandleHealth_Get_ReturnsOKStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handleHealth().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var got healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Status != "ok" {
		t.Errorf("status field = %q, want %q", got.Status, "ok")
	}
}

func Test_NewServer_UnknownPath_Returns404(t *testing.T) {
	srv := NewServer(":0", discardLogger(), BuildInfo{})

	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func Test_NewServer_WrongMethod_Returns405(t *testing.T) {
	srv := NewServer(":0", discardLogger(), BuildInfo{})

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func Test_NewServer_TimeoutsConfigured(t *testing.T) {
	srv := NewServer(":8080", discardLogger(), BuildInfo{})

	if srv.ReadHeaderTimeout == 0 {
		t.Error("ReadHeaderTimeout is unset; slow-header requests would hold connections open")
	}
	if srv.WriteTimeout == 0 {
		t.Error("WriteTimeout is unset")
	}
	if srv.MaxHeaderBytes == 0 {
		t.Error("MaxHeaderBytes is unset")
	}
}
