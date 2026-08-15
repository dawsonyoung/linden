# API Surface

> **Status:** Stub.

The client-facing HTTP contract, described for someone building against Linden.

Exact schemas, event names, and field types are normative in `docs/architecture/layer-interface-spec.md`. This page describes what each endpoint is for and what it guarantees; it links rather than restating, because a duplicated schema drifts.

## Stability

TODO: State the versioning and compatibility commitment before the first client ships.

## Endpoints

TODO: Specify as each endpoint lands. For each: purpose, request shape, response shape, error cases, and streaming behavior where applicable.

| Endpoint | Purpose | State |
|----------|---------|-------|
| `GET /health` | Liveness | Available |
| `GET /version` | Build identification | Available |
| `POST /chat` | Send a message, receive a reply | Planned — Stage A.7 |
| `GET /models` | List available models | Planned — Stage A.7 |

### `GET /health`

Reports whether the server is running and able to serve requests. Returns `200`
with a status field. Used by container orchestration and by clients deciding
whether the server is reachable.

It reports process liveness only. It does not indicate whether a model backend
is available.

### `GET /version`

Reports the running build: version, commit, and Go runtime version. Present so a
user can tell which build they are running when reporting a problem.

### Request identity

Every response carries an `X-Request-ID` header. A client may supply one to
correlate its own logs; if the supplied value is absent or unsafe to record, the
server substitutes a generated identifier.

## Streaming

Replies stream over Server-Sent Events. TODO: Describe the observable behavior a client can rely on — ordering, termination, and what happens if the connection drops. Link to the event contract in `docs/architecture/layer-interface-spec.md` for exact event names and payload shapes.

## Errors

TODO: Publish the client-facing error code list and what each means for the user. Internal detail is never included in a response body.
