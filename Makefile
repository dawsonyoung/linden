.PHONY: build build-web build-server dev test test-web validate validate-contracts validate-integration lint lint-web clean docs docs-check docs-serve docker-build docker-smoke

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT)

# POSIX only. Windows contributors run these inside the devcontainer.
# See docs/adr/0003-target-platform-policy.md.
#
# Targets whose component does not exist yet skip with a notice naming the
# stage that introduces it, rather than failing on a missing file.

# Build
build: build-web build-server

build-web:
	@if [ -f src/web/package.json ]; then \
	  cd src/web && npm ci && npm run build; \
	else \
	  echo "skip build-web: src/web/package.json not present (Stage C.3)"; \
	fi

build-server: build-web
	cd src && go build -ldflags "$(LDFLAGS)" -o ../bin/linden ./cmd/

# Development
dev: build-web
	cd src && go run ./cmd/

# Unit tests
test:
	cd src && go test -race ./...
	@$(MAKE) --no-print-directory test-web

test-web:
	@if [ -f src/web/package.json ]; then \
	  cd src/web && npm run check; \
	else \
	  echo "skip test-web: src/web/package.json not present (Stage C.3)"; \
	fi

# Validation suite
# validation/ is a separate module from src/; these run from validation/, not
# from src/ with a relative path.
validate: validate-contracts validate-integration

validate-contracts:
	@if [ -f validation/go.mod ]; then \
	  cd validation && go test -tags=contracts -race ./contracts/...; \
	else \
	  echo "skip validate-contracts: validation/go.mod not present (Stage A.1)"; \
	fi

validate-integration:
	@if [ -f validation/go.mod ]; then \
	  cd validation && go test -tags=integration -race ./integration/...; \
	else \
	  echo "skip validate-integration: validation/go.mod not present (Stage B.2)"; \
	fi

# Lint
# The lint job in .github/workflows/ci.yml invokes this target, so local and CI
# run the same checks by construction.
lint:
	cd src && go vet ./...
	@if [ -f validation/go.mod ]; then \
	  cd validation && go vet -tags=contracts ./...; \
	else \
	  echo "skip validation vet: validation/go.mod not present (Stage A.1)"; \
	fi
	@unformatted=$$(gofmt -l src validation); \
	if [ -n "$$unformatted" ]; then \
	  echo "gofmt reported unformatted files:"; \
	  echo "$$unformatted"; \
	  exit 1; \
	fi
	@$(MAKE) --no-print-directory lint-web

lint-web:
	@if [ -f src/web/package.json ]; then \
	  cd src/web && npm run check; \
	else \
	  echo "skip lint-web: src/web/package.json not present (Stage C.3)"; \
	fi

# Documentation
# docs-check catches what mdBook does not: pages absent from SUMMARY.md are
# silently unpublished, and relative links are never verified.
docs: docs-check verify-interfaces
	mdbook build

verify-interfaces:
	go run tools/docgen/main.go -verify

docs-check:
	@sh scripts/check-docs.sh

docs-serve: docs-check
	mdbook serve --open

# Container
check-docker:
	@docker info >/dev/null 2>&1 || (echo "error: docker daemon is not running"; exit 1)

docker-build: check-docker
	docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) -t linden:dev .

docker-smoke: docker-build
	@sh validation/smoke/smoke_test.sh

# Clean
clean:
	rm -rf bin/ src/web/build/ docs/.site/
