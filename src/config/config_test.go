package config

import (
	"log/slog"
	"testing"
)

func clearEnv(t *testing.T) {
	t.Helper()
	t.Setenv("LINDEN_ADDR", "")
	t.Setenv("LINDEN_LOG_LEVEL", "")
	t.Setenv("OLLAMA_URL", "")
}

func Test_Load_EnvUnset_AppliesDefaults(t *testing.T) {
	clearEnv(t)

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if got.Addr != defaultAddr {
		t.Errorf("Addr = %q, want %q", got.Addr, defaultAddr)
	}
	if got.LogLevel != slog.LevelInfo {
		t.Errorf("LogLevel = %v, want %v", got.LogLevel, slog.LevelInfo)
	}
	if got.OllamaURL != defaultOllamaURL {
		t.Errorf("OllamaURL = %q, want %q", got.OllamaURL, defaultOllamaURL)
	}
}

func Test_Load_EnvSet_OverridesDefaults(t *testing.T) {
	clearEnv(t)
	t.Setenv("LINDEN_ADDR", "127.0.0.1:9090")
	t.Setenv("LINDEN_LOG_LEVEL", "debug")
	t.Setenv("OLLAMA_URL", "https://ollama.internal:1234")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if got.Addr != "127.0.0.1:9090" {
		t.Errorf("Addr = %q, want %q", got.Addr, "127.0.0.1:9090")
	}
	if got.LogLevel != slog.LevelDebug {
		t.Errorf("LogLevel = %v, want %v", got.LogLevel, slog.LevelDebug)
	}
	if got.OllamaURL != "https://ollama.internal:1234" {
		t.Errorf("OllamaURL = %q, want %q", got.OllamaURL, "https://ollama.internal:1234")
	}
}

func Test_Load_WhitespaceValue_FallsBackToDefault(t *testing.T) {
	clearEnv(t)
	t.Setenv("LINDEN_ADDR", "   ")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if got.Addr != defaultAddr {
		t.Errorf("Addr = %q, want %q", got.Addr, defaultAddr)
	}
}

func Test_Load_InvalidValue_ReturnsError(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"address without port", "LINDEN_ADDR", "localhost"},
		{"address with too many colons", "LINDEN_ADDR", "a:b:c"},
		{"unknown log level", "LINDEN_LOG_LEVEL", "verbose"},
		{"url without scheme", "OLLAMA_URL", "localhost:11434"},
		{"url with unsupported scheme", "OLLAMA_URL", "ftp://localhost"},
		{"url without host", "OLLAMA_URL", "http://"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			t.Setenv(tt.key, tt.value)

			if _, err := Load(); err == nil {
				t.Fatalf("Load() error = nil, want error for %s=%q", tt.key, tt.value)
			}
		})
	}
}

func Test_ParseLevel_AllSupportedNames_Parse(t *testing.T) {
	tests := []struct {
		in   string
		want slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"warning", slog.LevelWarn},
		{"Error", slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := parseLevel(tt.in)
			if err != nil {
				t.Fatalf("parseLevel(%q) error = %v, want nil", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
