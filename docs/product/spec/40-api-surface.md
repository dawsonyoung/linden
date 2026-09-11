# API Surface

The client-facing HTTP contract for Linden.

Detailed Go interface shapes and internal error mappings are normative in `docs/architecture/layer-interface-spec.md`. This specification defines the client-facing HTTP routes, request/response formats, and streaming behaviors.

## Endpoints

| Endpoint | Purpose | Format | State |
|----------|---------|--------|-------|
| `GET /health` | Service liveness probe | JSON | Available |
| `GET /version` | Build version & commit | JSON | Available |
| `GET /models` | List locally available models | JSON | Available |
| `POST /chat` | Native Linden chat streaming | JSON / SSE | Available |
| `POST /v1/chat/completions` | OpenAI-compatible chat completion | JSON / SSE | Available |

---

### `GET /health`

Reports whether the server process is alive and responsive.
* **Response Status:** `200 OK`
* **Response Body:**
  ```json
  {"status":"ok"}
  ```

---

### `GET /version`

Reports the running application version, Git commit hash, and Go runtime version.
* **Response Status:** `200 OK`
* **Response Body:**
  ```json
  {
    "version": "v0.1.0",
    "commit": "0cc737b",
    "go_version": "go1.22.x"
  }
  ```

---

### `GET /models`

Queries the underlying inference engine (Ollama) and returns all currently installed models available for chat.
* **Response Status:** `200 OK`
* **Response Body:**
  ```json
  {
    "models": [
      {"name": "tinyllama:latest"},
      {"name": "qwen2.5:7b"}
    ]
  }
  ```

---

### `POST /chat`

The primary conversational endpoint used by the Linden web interface.
* **Headers:** `Content-Type: application/json`
* **Request Body:**
  ```json
  {
    "model": "tinyllama:latest",
    "messages": [
      {"role": "user", "content": "Hello Linden"}
    ]
  }
  ```
* **Response Headers:** `Content-Type: text/event-stream`
* **Streaming Event Format:**
  Incremental text tokens stream over Server-Sent Events (SSE). Each event carries a JSON chunk:
  ```
  data: {"text":"Hello","done":false}

  data: {"text":"! How","done":false}

  data: {"text":" can I help?","done":true}
  ```

---

### `POST /v1/chat/completions`

OpenAI-compatible translation endpoint supporting drop-in integration with external tools, IDE extensions, and agent frameworks.
* **Headers:** `Content-Type: application/json`
* **Request Body:**
  ```json
  {
    "model": "tinyllama:latest",
    "messages": [
      {"role": "user", "content": "Hello Linden"}
    ],
    "stream": true
  }
  ```
* **Response Headers:** `Content-Type: text/event-stream`
* **Streaming Event Format:**
  Conforms to OpenAI chunk streaming conventions:
  ```
  data: {"id":"chatcmpl-...","object":"chat.completion.chunk","created":1726056000,"model":"tinyllama:latest","choices":[{"index":0,"delta":{"content":"Hello"}}]}

  data: [DONE]
  ```

---

## Request Correlation

Every response carries an `X-Request-ID` header. Clients can provide their own tracking header in requests, or Linden generates a cryptographically random identifier.

## Client Error Handling

Client errors return standardized HTTP status codes (`400 Bad Request`, `404 Not Found`, `500 Internal Server Error`, `502 Bad Gateway` if inference is unreachable) with sanitized JSON payloads:
```json
{
  "code": "invalid_argument",
  "message": "model and messages are required"
}
```
Internal error traces are never exposed in API responses.
