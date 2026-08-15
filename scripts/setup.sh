#!/usr/bin/env sh
# Linden development environment check (Linux and macOS).
#
# Reports what is missing rather than installing it. Package managers vary and
# silent privileged installs are worse than a clear instruction.
#
# Usage: ./scripts/setup.sh [--pull-model]

set -eu

GO_MIN_MAJOR=1
GO_MIN_MINOR=22
NODE_MIN_MAJOR=20
MODEL=tinyllama

missing=0

step()  { printf '\n==> %s\n' "$1"; }
ok()    { printf '    OK: %s\n' "$1"; }
warn()  { printf '    MISSING: %s\n' "$1"; missing=$((missing + 1)); }
hint()  { printf '      %s\n' "$1"; }

step "Go >= ${GO_MIN_MAJOR}.${GO_MIN_MINOR}"
if command -v go >/dev/null 2>&1; then
  version=$(go env GOVERSION | sed 's/^go//')
  major=$(echo "$version" | cut -d. -f1)
  minor=$(echo "$version" | cut -d. -f2)
  if [ "$major" -gt "$GO_MIN_MAJOR" ] || { [ "$major" -eq "$GO_MIN_MAJOR" ] && [ "$minor" -ge "$GO_MIN_MINOR" ]; }; then
    ok "go ${version}"
  else
    warn "go ${version} is older than ${GO_MIN_MAJOR}.${GO_MIN_MINOR}"
    hint "https://go.dev/dl/"
  fi
else
  warn "go not found"
  hint "https://go.dev/dl/"
fi

step "Node >= ${NODE_MIN_MAJOR}"
if command -v node >/dev/null 2>&1; then
  version=$(node --version | sed 's/^v//')
  major=$(echo "$version" | cut -d. -f1)
  if [ "$major" -ge "$NODE_MIN_MAJOR" ]; then
    ok "node ${version}"
  else
    warn "node ${version} is older than ${NODE_MIN_MAJOR}"
    hint "https://nodejs.org/en/download"
  fi
else
  warn "node not found"
  hint "https://nodejs.org/en/download"
fi

step "Docker"
if command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1; then
  ok "docker daemon reachable"
else
  warn "docker not available"
  hint "Required for the runtime gate: make docker-smoke"
  hint "https://docs.docker.com/engine/install/"
fi

step "Ollama"
if command -v ollama >/dev/null 2>&1; then
  ok "$(ollama --version 2>&1 | head -1)"
  if [ "${1:-}" = "--pull-model" ]; then
    if ollama list 2>/dev/null | grep -q "$MODEL"; then
      ok "${MODEL} already present"
    else
      step "Pulling ${MODEL}"
      ollama pull "$MODEL"
      ok "${MODEL} pulled"
    fi
  else
    ollama list 2>/dev/null | grep -q "$MODEL" \
      && ok "${MODEL} present" \
      || hint "run with --pull-model to fetch ${MODEL}"
  fi
else
  warn "ollama not found"
  hint "Not required until Stage A.3"
  hint "https://ollama.com/download"
fi

printf '\n'
if [ "$missing" -eq 0 ]; then
  printf 'Environment ready.\n'
else
  printf '%s item(s) need attention. See the hints above.\n' "$missing"
fi
printf 'Gates: make lint && make test && make docker-smoke\n\n'
