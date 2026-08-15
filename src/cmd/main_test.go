package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func serverPort(t *testing.T, ts *httptest.Server) string {
	t.Helper()
	_, port, err := net.SplitHostPort(strings.TrimPrefix(ts.URL, "http://"))
	if err != nil {
		t.Fatalf("split test server URL %q: %v", ts.URL, err)
	}
	return port
}

func Test_ProbeHealthAt_ServerHealthy_ReturnsNil(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	if err := probeHealthAt("127.0.0.1:" + serverPort(t, ts)); err != nil {
		t.Fatalf("probeHealthAt() error = %v, want nil", err)
	}
}

func Test_ProbeHealthAt_WildcardAddress_ProbesLoopback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	if err := probeHealthAt(":" + serverPort(t, ts)); err != nil {
		t.Fatalf("probeHealthAt() error = %v, want nil", err)
	}
}

func Test_ProbeHealthAt_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	err := probeHealthAt("127.0.0.1:" + serverPort(t, ts))
	if err == nil {
		t.Fatal("probeHealthAt() error = nil, want error for 503")
	}
}

func Test_ProbeHealthAt_ServerDown_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	addr := "127.0.0.1:" + serverPort(t, ts)
	ts.Close()

	if err := probeHealthAt(addr); err == nil {
		t.Fatal("probeHealthAt() error = nil, want error when nothing is listening")
	}
}

func Test_ProbeHealthAt_MalformedAddress_ReturnsError(t *testing.T) {
	if err := probeHealthAt("not-an-address"); err == nil {
		t.Fatal("probeHealthAt() error = nil, want error for malformed address")
	}
}
