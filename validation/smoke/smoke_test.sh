#!/bin/sh
set -e

if ! docker info >/dev/null 2>&1; then
    echo "Error: Docker daemon is not running. Please start Docker and try again."
    exit 1
fi

# Find LAN IP to strictly verify 0.0.0.0 binding
LAN_IP=""
if command -v ip >/dev/null 2>&1; then
    LAN_IP=$(ip -4 route get 8.8.8.8 | awk '{print $7}' | tr -d '\n')
elif command -v ipconfig >/dev/null 2>&1; then
    LAN_IP=$(ipconfig | awk '/IPv4/ {print $NF}' | head -n 1 | tr -d '\r\n')
elif command -v ifconfig >/dev/null 2>&1; then
    LAN_IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -n 1)
fi

if [ -z "$LAN_IP" ]; then
    echo "Could not dynamically resolve LAN IP. Falling back to localhost."
    LAN_IP="localhost"
fi

echo "Starting dummy Ollama server for startup validation..."
cat << 'EOF' > /tmp/mock_ollama.go
package main
import (
	"fmt"
	"net/http"
)
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"models":[]}`)
	})
	http.ListenAndServe(":11434", nil)
}
EOF
go run /tmp/mock_ollama.go >/dev/null 2>&1 &
MOCK_PID=$!

echo "Starting smoke test container..."
# Ensure any old container is removed
docker rm -f linden-smoke >/dev/null 2>&1 || true

# Run container mapping port 8080 and setting OLLAMA_URL to the LAN IP mock
docker run -d --name linden-smoke -e OLLAMA_URL="http://${LAN_IP}:11434" -p 8080:8080 linden:dev

# Ensure cleanup on exit
trap 'echo "Cleaning up..." && kill $MOCK_PID 2>/dev/null && docker rm -f linden-smoke >/dev/null 2>&1' EXIT

echo "Waiting for /health endpoint..."
MAX_TRIES=15
TRIES=0
while [ $TRIES -lt $MAX_TRIES ]; do
    if curl -s http://localhost:8080/health >/dev/null 2>&1; then
        break
    fi
    sleep 1
    TRIES=$((TRIES + 1))
done

if [ $TRIES -eq $MAX_TRIES ]; then
    echo "Container failed to become healthy."
    docker logs linden-smoke
    exit 1
fi

echo "Container is healthy. Running integration suite against LAN IP: $LAN_IP"
export LINDEN_TEST_URL="http://${LAN_IP}:8080"

# Run the integration suite against the LAN URL
cd ../..
go test -v -tags=integration ./validation/integration/...

echo "Smoke test passed."
