// Package api exposes the HTTP surface. It depends only on interfaces from
// orchestrator; health and version are self-contained and have no dependency.
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/dawsonyoung/linden/orchestrator"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	// TODO: SSE responses outlive this. Streaming routes need per-route write
	// deadlines when chat lands in Stage A.7.
	writeTimeout   = 30 * time.Second
	idleTimeout    = 60 * time.Second
	maxHeaderBytes = 1 << 20
)

// BuildInfo identifies the running binary. Values are injected at build time.
type BuildInfo struct {
	Version   string
	Commit    string
	GoVersion string
}

// NewServer wires routes and middleware and applies timeouts. The caller owns
// starting and stopping the returned server.
func NewServer(addr string, logger *slog.Logger, build BuildInfo, chatService orchestrator.ChatService) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth())
	mux.HandleFunc("GET /version", handleVersion(build))

	return &http.Server{
		Addr:              addr,
		Handler:           requestID(accessLog(logger, mux)),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
