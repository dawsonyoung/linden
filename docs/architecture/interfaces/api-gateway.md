# API Gateway Technical Specification

## 1. Overview & Boundary

The `api` layer (`src/api/`) serves as Linden's external HTTP gateway. It parses incoming client HTTP requests, validates payload sizes and schema limits, translates requests into internal `orchestrator.ChatService` calls, and writes JSON or Server-Sent Events (SSE) responses.

- **Package:** `github.com/dawsonyoung/linden/api`
- **Dependency:** Depends only on `orchestrator` and `errs`. No direct imports of `inference`, `storage`, or `cmd`.
- **Consumers:** SvelteKit frontend (`src/web/`), LAN client devices, external LLM tooling.

---

## 2. HTTP Endpoint Routing Table

| Method | Path | Request Content-Type | Response Content-Type | Handler Function | Purpose |
|---|---|---|---|---|---|
| `GET` | `/health` | None | `application/json` | `handleHealth()` | Liveness check returning `{"status":"ok"}` |
| `GET` | `/version` | None | `application/json` | `handleVersion()` | Returns Git commit, version, and build info |
| `GET` | `/models` | None | `application/json` | `handleModels(svc)` | Lists locally installed models available for inference |
| `POST` | `/chat` | `application/json` | `text/event-stream` / `application/json` | `handleChat(svc)` | Native streaming multi-turn chat generation |
| `POST` | `/v1/chat/completions` | `application/json` | `text/event-stream` / `application/json` | `handleOpenAICompat(svc)` | OpenAI-compatible chat completion endpoint |

---

## 3. Request & Response Schemas

### Request Headers
| Header | Required | Valid Values | Behavior |
|---|---|---|---|
| `Content-Type` | Required for `POST` | `application/json` | Requests with missing/invalid Content-Type return `400 Bad Request`. |
| `X-Request-ID` | Optional | String (max 64 chars) | Propagated to response headers and server log contexts for traceability. |

### Native Chat Request Schema (`POST /chat`)
| Field | Type | Required | Constraints | Description |
|---|---|---|---|---|
| `model` | `string` | Yes | 1–100 chars | Name of local model (e.g. `"tinyllama"`, `"llama3"`). |
| `messages` | `[]chatMessage` | Yes | 1–1,000 items | Chronological list of conversational turns. |
| `messages[].role` | `string` | Yes | `"system"`, `"user"`, `"assistant"` | Role of the message producer. |
| `messages[].content`| `string` | Yes | Non-empty string | Text content of the conversational turn. |
| `stream` | `boolean` | No | Default `true` | When true, streams via SSE; when false, returns terminal JSON. |
| `sessionId` | `string` | No | Valid UUID / hex | Optional session identifier for history association. |

### Native Chat Chunk Response (`text/event-stream`)
| Event Name | Data Payload Schema | Terminal | Description |
|---|---|---|---|
| *default* (`data:`) | `{"text": "string"}` | No | Incremental token delivered during generation. |
| `event: error` | `{"error": "string"}` | Yes | Emitted when generation aborts due to an unrecoverable error. |
| `event: done` | `[DONE]` | Yes | Sent to signal normal completion of token stream. |

### Models Response Schema (`GET /models`)
```json
{
  "models": [
    { "name": "tinyllama" },
    { "name": "llama3:latest" }
  ]
}
```

### Error Response Schema (`application/json`)
```json
{
  "error": "human-readable safe message",
  "code": "invalid_argument"
}
```

---

## 4. Execution Semantics & Concurrency Rules

1. **Payload Body Limit:** All request bodies are strictly constrained to a maximum of **1 MB** via `http.MaxBytesReader`. Payloads exceeding 1 MB are rejected with `400 Bad Request` before decoding.
2. **Context Propagation & Cancellation:** Every request passes `r.Context()` to the orchestrator. If the client drops the TCP connection or cancels the request, the context cancellation cancels backend inference immediately.
3. **SSE Flushing:** Handlers assert `http.Flusher`. If unsupported, the server responds with `500 Internal Server Error`. Each token chunk is flushed immediately upon receipt from `onChunk`.
4. **Zero-Logging Privacy:** Prompt messages, generated tokens, and user inputs are strictly excluded from logging at all log levels.

---

## 5. Error Taxonomy & HTTP Status Code Mapping

The API layer translates internal `errs.Code` classifications to standard HTTP status codes:

| `errs.Code` | HTTP Status | Returned `code` Field | Scenario Trigger |
|---|---|---|---|
| `errs.InvalidArgument` | `400 Bad Request` | `"invalid_argument"` | Empty messages, missing model, malformed JSON, payload > 1MB. |
| `errs.NotFound` | `404 Not Found` | `"not_found"` | Requested model is not installed or session does not exist. |
| `errs.Unavailable` | `503 Service Unavailable` | `"unavailable"` | Inference daemon (Ollama) is offline or refusing connections. |
| `errs.DeadlineExceeded`| `503 Service Unavailable` | `"deadline_exceeded"` | Inference backend timed out fulfilling request. |
| `errs.Internal` (or unclassified) | `500 Internal Server Error` | `"internal"` | Unexpected server-side fault or unhandled panic. |

---

## 6. Contract Test Verification Matrix

Every behavior in this specification is verified in `validation/contracts/`:

| Test Function | Contract Verified |
|---|---|
| `Test_API_Health` | Validates `GET /health` returns `200 OK` and `{"status":"ok"}`. |
| `Test_API_Version` | Validates `GET /version` returns JSON with version fields. |
| `Test_API_Models` | Validates `GET /models` lists models and handles orchestrator errors. |
| `Test_API_Chat_Stream` | Validates `POST /chat` streams SSE chunks and terminates with `[DONE]`. |
| `Test_API_Chat_EmptyMessages` | Validates `POST /chat` rejects empty message list with `400` and `invalid_argument`. |
| `Test_API_Chat_MissingModel` | Validates `POST /chat` rejects empty model string with `400`. |
| `Test_API_OpenAI_Compat` | Validates `POST /v1/chat/completions` formats SSE chunks in OpenAI delta format. |
