// Command linden runs the Linden HTTP server.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dawsonyoung/linden/api"
	"github.com/dawsonyoung/linden/config"
	"github.com/dawsonyoung/linden/inference"
	"github.com/dawsonyoung/linden/orchestrator"
	"github.com/dawsonyoung/linden/storage"
)

// Injected via -ldflags at build time; empty values are normalized by the api layer.
var (
	version string
	commit  string
)

const (
	shutdownTimeout    = 10 * time.Second
	healthProbeTimeout = 2 * time.Second
)

func main() {
	health := flag.Bool("health", false, "probe the local health endpoint and exit; used by the container HEALTHCHECK")
	flag.Parse()

	if *health {
		if err := probeHealth(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	if err := run(); err != nil {
		slog.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func probeHealth() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return probeHealthAt(cfg.Addr)
}

// probeHealthAt reports whether the server at addr is serving health checks.
// A wildcard listen address is probed over loopback.
func probeHealthAt(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("parse address: %w", err)
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}

	client := &http.Client{Timeout: healthProbeTimeout}
	resp, err := client.Get("http://" + net.JoinHostPort(host, port) + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health endpoint returned %d", resp.StatusCode)
	}
	return nil
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	inf, err := inference.NewOllama(inference.OllamaConfig{
		BaseURL:         cfg.OllamaURL,
		ResponseTimeout: 30 * time.Second, // Default to a reasonable headers wait
	})
	if err != nil {
		return fmt.Errorf("configure inference: %w", err)
	}

	store, err := storage.NewFileStore(cfg.DataDir)
	if err != nil {
		return fmt.Errorf("configure storage: %w", err)
	}

	chatService := orchestrator.NewService(inf, store)

	// Validate inference engine connection before serving
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if _, err := chatService.ListModels(pingCtx); err != nil {
		return fmt.Errorf("fatal: inference engine unreachable at startup: %w", err)
	}

	srv := api.NewServer(cfg.Addr, logger, api.BuildInfo{Version: version, Commit: commit}, chatService)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		logger.Info("server starting", slog.String("addr", cfg.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	logger.Info("server stopped")
	return nil
}
