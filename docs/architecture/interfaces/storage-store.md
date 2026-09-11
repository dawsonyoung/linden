# Storage Store Technical Specification

## 1. Overview & Boundary

The `storage` layer (`src/storage/`) manages persistence of user sessions, conversation turns, and local configuration. It guarantees strict isolation, transparent local file storage, and deterministic data purging.

- **Package:** `github.com/dawsonyoung/linden/storage`
- **Exported Interface:** `Store`
- **Dependencies:** Leaf layer. Depends only on `errs` and Go standard library. Never imports `api`, `orchestrator`, or `inference`.
- **Consumers:** Consumed by `src/orchestrator/`.

---

## 2. Interface Definition

```go
type Store interface {
    // SaveTurn persists a single turn for the given session ID.
    SaveTurn(ctx context.Context, sessionID string, turn Turn) error

    // LoadSession retrieves all turns for a specific session chronologically.
    // Returns errs.NotFound if the session does not exist.
    LoadSession(ctx context.Context, sessionID string) ([]Turn, error)

    // ListSessions returns a chronologically ordered list of all known sessions.
    ListSessions(ctx context.Context) ([]Session, error)
}
```

---

## 3. Data Structures & Schemas

### `storage.Session`
| Field | Type | Description |
|---|---|---|
| `ID` | `string` | Unique session identifier (UUID / alphanumeric hex). |
| `CreatedAt` | `time.Time` | Timestamp when the session was created (UTC RFC3339). |
| `UpdatedAt` | `time.Time` | Timestamp of the most recent turn in the session (UTC RFC3339). |

### `storage.Turn`
| Field | Type | Required | Constraints | Description |
|---|---|---|---|---|
| `Role` | `string` | Yes | `"user"` or `"assistant"` | Role of the entity producing this turn. |
| `Content` | `string` | Yes | Non-empty | Text content of the conversation turn. |
| `Timestamp` | `time.Time` | Yes | Valid UTC time | Exact time the turn was received or completed. |

---

## 4. Method Specifications

### `SaveTurn(ctx context.Context, sessionID string, turn Turn) error`
- **Preconditions:** `sessionID` must not be empty. `turn.Content` must not be empty.
- **Postconditions:** Appends the turn to the session record. If the session does not exist, creates it atomically. Updates session `UpdatedAt` timestamp.
- **Return Values:** Nil on success; wrapped `errs.Error` on failure.

### `LoadSession(ctx context.Context, sessionID string) ([]Turn, error)`
- **Preconditions:** `sessionID` must not be empty.
- **Postconditions:** Returns all turns ordered chronologically. Returns empty slice or `errs.NotFound` if session does not exist.
- **Return Values:**
  - `[]Turn`: Chronological turns.
  - `error`: Nil on success; `errs.NotFound` if session ID is missing from disk.

### `ListSessions(ctx context.Context) ([]Session, error)`
- **Preconditions:** `ctx` is valid.
- **Postconditions:** Returns metadata for all existing sessions, sorted newest first (`UpdatedAt` descending). Returns empty slice if no sessions exist.
- **Return Values:**
  - `[]Session`: Session summaries.
  - `error`: Nil on success; wrapped `errs.Error` on disk failure.

---

## 5. Storage Formats & Privacy Guarantees

1. **Local Filesystem Placement:** Sessions are stored under `data/sessions/` on the local machine using atomic write operations (`.tmp` write followed by `os.Rename`).
2. **Deterministic Data Purge:** When a session or turn is deleted, the file is unlinked immediately (`os.Remove`). No orphaned chunks or tombstone records remain on disk.
3. **No Unencrypted Cloud Sync:** The storage directory is local to the host and excluded from cloud syncing mechanisms.

---

## 6. Error Taxonomy & Mapping Table

| Scenario | `errs.Code` | Notes |
|---|---|---|
| Empty `sessionID` passed | `errs.InvalidArgument` | Validated at boundary before I/O. |
| Empty `turn.Content` passed | `errs.InvalidArgument` | Empty turns are rejected. |
| Requested session does not exist | `errs.NotFound` | Returned by `LoadSession`. |
| Disk full or permission error | `errs.ResourceExhausted` / `errs.Internal` | I/O failure mapped safely. |

---

## 7. Contract Test Verification Matrix

Every requirement in this specification is verified in `validation/contracts/`:

| Test Function | Contract Verified |
|---|---|
| `Test_Storage_SaveAndLoadSession` | Verifies saving turns and retrieving in chronological order. |
| `Test_Storage_LoadNonexistent_ReturnsNotFound` | Verifies `LoadSession` on missing session yields `errs.NotFound`. |
| `Test_Storage_ListSessions_Order` | Verifies sessions are returned sorted by `UpdatedAt` descending. |
| `Test_Storage_EmptyInput_Validation` | Verifies empty session ID or turn rejects with `errs.InvalidArgument`. |
