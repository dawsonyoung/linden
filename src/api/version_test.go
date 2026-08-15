package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

func Test_HandleVersion_BuildInfoSet_ReturnsProvidedValues(t *testing.T) {
	build := BuildInfo{Version: "1.2.3", Commit: "abc1234", GoVersion: "go1.22.0"}

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	handleVersion(build).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got versionResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Version != "1.2.3" {
		t.Errorf("version = %q, want %q", got.Version, "1.2.3")
	}
	if got.Commit != "abc1234" {
		t.Errorf("commit = %q, want %q", got.Commit, "abc1234")
	}
	if got.GoVersion != "go1.22.0" {
		t.Errorf("go_version = %q, want %q", got.GoVersion, "go1.22.0")
	}
}

func Test_HandleVersion_BuildInfoUnset_ReturnsDevDefaults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	handleVersion(BuildInfo{}).ServeHTTP(rec, req)

	var got versionResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Version != "dev" {
		t.Errorf("version = %q, want %q", got.Version, "dev")
	}
	if got.Commit != "unknown" {
		t.Errorf("commit = %q, want %q", got.Commit, "unknown")
	}
	if got.GoVersion != runtime.Version() {
		t.Errorf("go_version = %q, want %q", got.GoVersion, runtime.Version())
	}
}
