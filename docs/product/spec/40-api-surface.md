# API Surface

> **Status:** Stub.

The client-facing HTTP contract. This page describes the surface a client can depend on. Layer-internal interfaces are in `docs/architecture/layer-interface-spec.md` and are not part of this contract.

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

Replies stream over Server-Sent Events. TODO: Restate the client-visible event contract here once `docs/architecture/layer-interface-spec.md` stabilizes, describing observable behavior only.

## Errors

TODO: Publish the client-facing error code list and what each means for the user. Internal detail is never included in a response body.
