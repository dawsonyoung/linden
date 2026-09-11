# Error Taxonomy Technical Specification

## 1. Overview & Architectural Charter

The `errs` package (`src/errs/`) is Linden's **shared error kernel**. It defines the closed taxonomy of error codes used across all architectural layers.

- **Package:** `github.com/dawsonyoung/linden/errs`
- **Dependencies:** None. Standard library only. Zero I/O, zero network or transport dependencies.
- **Architectural Rule:** Any layer (including leaf layers) may import `src/errs`. No layer imports another layer for error definitions.
- **Governing Decision:** [ADR-0004: Shared Error Taxonomy](../../adr/0004-shared-error-taxonomy.md).

---

## 2. Core Go Type Definitions

```go
type Code string

const (
    InvalidArgument   Code = "invalid_argument"
    Unauthenticated   Code = "unauthenticated"
    PermissionDenied  Code = "permission_denied"
    NotFound          Code = "not_found"
    Conflict          Code = "conflict"
    ResourceExhausted Code = "resource_exhausted"
    Unavailable       Code = "unavailable"
    DeadlineExceeded  Code = "deadline_exceeded"
    Internal          Code = "internal"
)

type Error struct {
    Code    Code
    Message string
    cause   error
}
```

---

## 3. Normative Taxonomy & Mapping Table

The set of codes is closed. Adding or removing a code requires updating this specification and `src/errs/errs.go` in the same commit.

| `Code` String | Semantic Meaning | Recommended HTTP Status | Client Visibility |
|---|---|---|---|
| `"invalid_argument"` | Client specified an invalid argument (empty messages, missing model, invalid schema). | `400 Bad Request` | Safe to return validation message. |
| `"unauthenticated"` | Request lacks valid authentication credentials. | `401 Unauthorized` | Safe generic message. |
| `"permission_denied"`| Caller lacks permission to execute the requested operation. | `403 Forbidden` | Safe generic message. |
| `"not_found"` | A specified resource or model was not found. | `404 Not Found` | Safe to name missing resource/model. |
| `"conflict"` | Concurrency conflict or resource state collision. | `409 Conflict` | Safe generic message. |
| `"resource_exhausted"` | Host memory, thread budget, or storage quota exceeded. | `429 Too Many Requests` / `507` | Safe backpressure notice. |
| `"unavailable"` | Service, inference engine, or dependency is offline/unreachable. | `503 Service Unavailable` | Actionable recovery message. |
| `"deadline_exceeded"` | Operation timed out before completing. | `504 Gateway Timeout` / `503` | Safe timeout message. |
| `"internal"` | Internal system defect, unexpected panic, or unclassified failure. | `500 Internal Server Error` | Generic message; internal cause hidden. |

---

## 4. Propagation & Wrapping Rules

1. **`errs.New(code, message)`:** Constructs an error with the specified code and diagnostic message.
2. **`errs.Wrap(code, message, err)`:** Annotates `err` with a new `code`, preserving the underlying error chain for `errors.Is` and `errors.As`. If `err` is nil, returns nil.
3. **`errs.CodeOf(err)`:** Evaluates the outermost classified error in the error chain. If `err` carries no taxonomy code, `CodeOf` reports `errs.Internal`. An unclassified error is treated as a server defect.
4. **`errs.Is(err, code)`:** Convenience function checking whether `CodeOf(err) == code`.
5. **No User Content in Messages:** The `Message` field reaches server operational logs. It must **never** contain prompt text, user tokens, or private user data.
