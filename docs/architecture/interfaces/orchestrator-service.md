# Orchestrator Service Technical Specification

## 1. Overview & Boundary

The `orchestrator` layer (`src/orchestrator/`) manages request routing, context assembly, policy enforcement, and interaction between the API gateway, local inference backend, and storage layer.

- **Package:** `github.com/dawsonyoung/linden/orchestrator`
- **Exported Interface:** `ChatService`
- **Dependencies:** Depends only on `inference.Client`, `storage.Store`, and `errs`. No imports of `api` or `cmd`.
- **Consumers:** Consumed exclusively by `src/api/` handlers.

---

## 2. Interface Definition

```go
type ChatService interface {
    // ListModels reports the models available for chat completion.
    ListModels(ctx context.Context) ([]Model, error)

    // ChatStream generates a reply, delivering incremental chunks to onChunk.
    ChatStream(ctx context.Context, req Request, onChunk func(Chunk) error) (Result, error)
}
```

---

## 3. Data Structures & Schemas

### `orchestrator.Request`
| Field | Type | Required | Constraints | Description |
|---|---|---|---|---|
| `Model` | `string` | Yes | Non-empty | Target model name (e.g. `"tinyllama"`). |
| `Messages` | `[]Message` | Yes | Minimum 1 item | Sequential conversation turns. |
| `SessionID` | `string` | No | Optional | Session ID for storage persistence. |

### `orchestrator.Message`
| Field | Type | Required | Valid Values | Description |
|---|---|---|---|---|
| `Role` | `Role` | Yes | `"system"`, `"user"`, `"assistant"` | Sender role in conversation. |
| `Content` | `string` | Yes | Non-empty | Content of the message. |

### `orchestrator.Chunk`
| Field | Type | Description |
|---|---|---|
| `Text` | `string` | Incremental token fragment emitted during generation. |

### `orchestrator.Result`
| Field | Type | Description |
|---|---|---|
| `Model` | `string` | The model that produced the generation. |
| `FinishReason` | `FinishReason` | Cause of completion: `"stop"` (normal) or `"canceled"` (aborted). |

---

## 4. Method Specifications

### `ListModels(ctx context.Context) ([]Model, error)`
- **Preconditions:** `ctx` is not nil.
- **Postconditions:** Returns slice of `Model` reported by the underlying inference client. If inference client fails, wraps and propagates error with taxonomy preserved.
- **Parameters:**
  - `ctx`: Standard execution context with timeout or cancellation.
- **Return Values:**
  - `[]Model`: List of available models. Returns empty slice (not nil) if no models installed.
  - `error`: Nil on success; error carrying `errs.Code` on failure.

### `ChatStream(ctx context.Context, req Request, onChunk func(Chunk) error) (Result, error)`
- **Preconditions:**
  - `req.Messages` must contain at least 1 message.
  - `req.Model` must not be empty.
  - `onChunk` callback must not be nil.
- **Postconditions:**
  - Emits tokens sequentially via `onChunk`.
  - Returns `Result` with `FinishReason == FinishStop` on successful completion.
  - Returns `Result{FinishReason: FinishCanceled}` when context is canceled.
- **Streaming Callback Contract:**
  1. `onChunk` is called synchronously on the calling goroutine, in order, and strictly before `ChatStream` returns. Consumers require no synchronization.
  2. If `onChunk` returns a non-nil error, `ChatStream` halts delivery immediately, terminates inference, and returns the unwrapped error to allow sentinel matching.
  3. The implementation must check `ctx.Err()` before each chunk delivery. Cancellation takes precedence over transport classification.

---

## 5. Error Taxonomy & Status Mapping

| Trigger Condition | `errs.Code` | Notes |
|---|---|---|
| `len(req.Messages) == 0` | `errs.InvalidArgument` | Validated before calling inference client. |
| `req.Model == ""` | `errs.InvalidArgument` | Validated before calling inference client. |
| Backend inference client unreachable | `errs.Unavailable` | Propagated from inference layer without alteration. |
| Context canceled (`ctx.Done()`) | `context.Canceled` | Returns `FinishCanceled` with `ctx.Err()`. |
| Unhandled internal error | `errs.Internal` | Sanitized to prevent internal details from leaking. |

---

## 6. Contract Test Verification Matrix

Every requirement in this specification is verified in `validation/contracts/`:

| Test Function | Contract Verified |
|---|---|
| `Test_Orchestrator_ListModels_Success` | Verifies models returned match underlying inference provider. |
| `Test_Orchestrator_ListModels_PropagatesError` | Verifies inference error taxonomy is preserved. |
| `Test_Orchestrator_ChatStream_Success` | Verifies chunks deliver in order and ends with `FinishStop`. |
| `Test_Orchestrator_ChatStream_EmptyMessages` | Verifies immediate rejection with `errs.InvalidArgument`. |
| `Test_Orchestrator_ChatStream_CallbackError_Aborts` | Verifies non-nil callback return halts generation and unwraps error. |
| `Test_Orchestrator_ChatStream_Cancellation` | Verifies `ctx.Cancel()` yields `FinishCanceled` and `context.Canceled`. |
