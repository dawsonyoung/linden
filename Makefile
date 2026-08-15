.PHONY: build build-web build-server dev test validate lint clean docs docs-serve

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT)

# Build everything
build: build-web build-server

build-web:
	cd src/web && npm ci && npm run build

build-server: build-web
	cd src && go build -ldflags "$(LDFLAGS)" -o ../bin/linden ./cmd/

# Development
dev: build-web
	cd src && OLLAMA_URL=http://localhost:11434 go run ./cmd/

# Unit tests
test:
	cd src && go test ./...
	cd src/web && npm run check

# Validation suite
validate: validate-contracts validate-integration validate-smoke

validate-contracts:
	cd src && go test -tags=contracts ../validation/contracts/...

validate-integration:
	cd src && go test -tags=integration ../validation/integration/...

validate-smoke:
	cd validation/smoke && ./run.sh

# Lint
lint:
	cd src && go vet ./...
	cd src/web && npm run check

# Documentation
docs:
	@echo "TODO: implement tools/docgen (Stage 0.6). See docs/adr/0001-documentation-publishing-toolchain.md"
	@exit 1

docs-serve:
	@echo "TODO: implement tools/docgen (Stage 0.6). See docs/adr/0001-documentation-publishing-toolchain.md"
	@exit 1

# Clean
clean:
	rm -rf bin/ src/web/build/ docs/.site/
