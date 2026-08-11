.PHONY: build build-web build-server dev test validate lint clean

# Build everything
build: build-web build-server

build-web:
	cd src/web && npm ci && npm run build

build-server: build-web
	cd src && go build -o ../bin/linden ./cmd/

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

# Clean
clean:
	rm -rf bin/ src/web/build/
