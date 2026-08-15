.PHONY: build build-web build-server dev test validate lint clean docs docs-serve docker-build docker-smoke

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

# Container
docker-build:
	docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) -t linden:dev .

docker-smoke: docker-build
	-docker rm -f linden-smoke
	docker run -d --name linden-smoke -p 8080:8080 linden:dev
	for i in $$(seq 1 20); do \
	  code=$$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8080/health || true); \
	  if [ "$$code" = "200" ]; then echo "health 200 after $${i}s"; break; fi; \
	  if [ $$i -eq 20 ]; then echo "health check failed"; docker logs linden-smoke; docker rm -f linden-smoke; exit 1; fi; \
	  sleep 1; \
	done
	docker stop -t 15 linden-smoke
	docker logs linden-smoke
	docker rm -f linden-smoke

# Clean
clean:
	rm -rf bin/ src/web/build/ docs/.site/
