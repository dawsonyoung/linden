# Inference Client Technical Specification

## 1. Overview & Boundary

The `inference` layer (`src/inference/`) abstracts local LLM backends (such as Ollama or local GGUF runtimes). It translates internal inference requests into provider-specific HTTP wire calls, consumes chunked responses, and normalizes transport errors into Linden's error taxonomy.

- **Package:** `github.com/dawsonyoung/linden/inference`
- **Exported Interface:** `Client`
- **Dependencies:** Leaf layer. Depends only on `errs`. Never imports `api`, `orchestrator`, or `storage`.
- **Consumers:** Consumed exclusively by `src/orchestrator/`.

---

## 2. Interface Definition

```go
type Client interface {
    // ListModels reports the models the provider can serve.
    ListModels(ctx context.Context) ([]Model, error)

    // ChatStream generates a reply, delivering it to onChunk as it arrives.
    ChatStream(ctx context.Context, req Request, onChunk func(Chunk) error) (Result, error)
}
```

---

## 3. Data Structures & Schemas

### `inference.Request`
| Field | Type | Required | Constraints | Description |
|---|---|---|---|---|
| `Model` | `string` | Yes | Non-empty | Target model name installed in local backend. |
| `Messages` | `[]Message` | Yes | Minimum 1 item | Context turns sent to model. |

### `inference.Message`
| Field | Type | Required | Valid Values | Description |
|---|---|---|---|---|
| `Role` | `Role` | Yes | `"system"`, `"user"`, `"assistant"` | Origin of the message. |
| `Content` | `string` | Yes | Non-empty | Text content of the prompt turn. |

### `inference.Chunk`
| Field | Type | Description |
|---|---|---|
| `Text` | `string` | Incremental token stream fragment emitted by the inference engine. |

### `inference.Result`
| Field | Type | Description |
|---|---|---|
| `Model` | `string` | The model that executed the completion. |
| `FinishReason` | `FinishReason` | `"stop"` (normal completion) or `"canceled"` (context canceled). |

---

## 4. Execution Semantics & Concurrency Rules

1. **Callback Delivery Order:** `onChunk` is called sequentially on the calling goroutine. Consumers require no mutexes or synchronization primitives.
2. **Consumer Abort Handling:** If `onChunk` returns an error, `ChatStream` aborts stream consumption immediately, closes provider connections, and returns the unwrapped error to allow sentinel matching.
3. **Cancellation Precedence:** The client checks `ctx.Err()` before each chunk delivery and before classifying any network or transport failure. Cancellation takes precedence over connection errors, returning `FinishCanceled` with `context.Canceled`.
4. **Zero Result on Error:** On error paths (other than context cancellation), `Result` returns the zero value (`Result{}`).
5. **Zero User-Content Logging:** Prompts, system messages, and generated response tokens must **never** be logged at any level.

---

## 5. Supported Ollama Container Topologies

| Topology Mode | Environment Variable | Target Endpoint | Description |
|---|---|---|---|
| **Host-Native Ollama** | `OLLAMA_URL=http://host.docker.internal:11434` | Host OS Daemon | Connects to native Ollama running on host for direct GPU hardware acceleration. |
| **Compose Sibling** | `OLLAMA_URL=http://ollama:11434` | Sibling Container | Connects to an Ollama container sharing the same Docker network bridge. |
| **Local Bare-Metal** | `OLLAMA_URL=http://127.0.0.1:11434` | Loopback Port | Connects to Ollama running as a standard local process during bare-metal execution. |

---

## 6. Provider Error to `errs.Code` Mapping Table

| Provider Transport Event | Internal HTTP Status | Mapped `errs.Code` | Notes |
|---|---|---|---|
| Request with 0 messages | N/A | `errs.InvalidArgument` | Rejected before making any HTTP request. |
| Model not found in backend | HTTP 404 | `errs.NotFound` | Indicates model must be pulled locally. |
| Provider offline / connection refused | TCP dial error | `errs.Unavailable` | Actionable message names backend failure. |
| Provider response timeout | HTTP client timeout | `errs.DeadlineExceeded`| Triggered if backend stalls generation. |
| Context canceled (`ctx.Done()`) | Context cancel | `context.Canceled` | Returns `FinishCanceled`. |
| Malformed JSON payload from backend | HTTP 200 invalid JSON | `errs.Internal` | Provider returned corrupt stream. |

---

## 7. Contract Test Verification Matrix

Every requirement in this specification is verified in `validation/contracts/`:

| Test Function | Contract Verified |
|---|---|
| `Test_Inference_ListModels_Success` | Verifies models returned match provider models. |
| `Test_Inference_ListModels_Unreachable` | Verifies connection refusal returns `errs.Unavailable`. |
| `Test_Inference_ChatStream_HappyPath` | Verifies chunk stream order and terminal `FinishStop`. |
| `Test_Inference_ChatStream_EmptyMessages` | Verifies early rejection with `errs.InvalidArgument`. |
| `Test_Inference_ChatStream_CallbackError_Aborts` | Verifies callback error stops stream and unwraps. |
| `Test_Inference_ChatStream_ContextCancellation` | Verifies cancellation returns `FinishCanceled` and `context.Canceled`. |
| `Test_Inference_ChatStream_ProviderTimeout` | Verifies timeout yields `errs.DeadlineExceeded`. |
