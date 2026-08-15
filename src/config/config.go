// Package config loads runtime configuration from the environment.
package config

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strings"
)

const (
	defaultAddr      = ":8080"
	defaultLogLevel  = "info"
	defaultOllamaURL = "http://localhost:11434"
)

// Config holds validated runtime settings.
type Config struct {
	Addr      string
	LogLevel  slog.Level
	OllamaURL string
}

// Load reads configuration from the environment, applying defaults for unset
// values. It returns an error naming the offending variable rather than
// falling back silently, so a misconfigured deployment fails at startup.
func Load() (Config, error) {
	addr := lookup("LINDEN_ADDR", defaultAddr)
	if err := validateAddr(addr); err != nil {
		return Config{}, fmt.Errorf("LINDEN_ADDR: %w", err)
	}

	level, err := parseLevel(lookup("LINDEN_LOG_LEVEL", defaultLogLevel))
	if err != nil {
		return Config{}, fmt.Errorf("LINDEN_LOG_LEVEL: %w", err)
	}

	ollamaURL := lookup("OLLAMA_URL", defaultOllamaURL)
	if err := validateURL(ollamaURL); err != nil {
		return Config{}, fmt.Errorf("OLLAMA_URL: %w", err)
	}

	return Config{Addr: addr, LogLevel: level, OllamaURL: ollamaURL}, nil
}

func lookup(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			return trimmed
		}
	}
	return fallback
}

func validateAddr(addr string) error {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		return fmt.Errorf("must be host:port: %w", err)
	}
	return nil
}

func parseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unknown level %q, want debug, info, warn, or error", s)
	}
}

func validateURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("not a valid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("scheme must be http or https, got %q", u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("missing host")
	}
	return nil
}
