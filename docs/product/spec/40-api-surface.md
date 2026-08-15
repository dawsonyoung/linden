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
| `GET /health` | Liveness | Planned — Stage 0.2 |
| `GET /version` | Build identification | Planned — Stage 0.2 |
| `POST /chat` | Send a message, receive a reply | Planned — Stage A.7 |
| `GET /models` | List available models | Planned — Stage A.7 |

## Streaming

Replies stream over Server-Sent Events. TODO: Describe the observable behavior a client can rely on — ordering, termination, and what happens if the connection drops. Link to the event contract in `docs/architecture/layer-interface-spec.md` for exact event names and payload shapes.

## Errors

TODO: Publish the client-facing error code list and what each means for the user. Internal detail is never included in a response body.
