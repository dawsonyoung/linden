# Capabilities

Each capability lists its originating requirement and its current observable state in the system.

| Capability | Origin | State |
|-----------|--------|-------|
| Local chat | PRD-0001 | Available |
| Model selection | PRD-0001 | Available |
| OpenAI API compatibility | Feature | Available |
| Zero-config LAN discovery | Feature | Available |
| Document retrieval (RAG) | Planned | Track C |
| Multi-device access control | Planned | Track E |

## Local chat

Users can carry out streaming multi-turn conversations with locally running language models.

- **Message Exchange:** Accepts user input, constructs conversation turns (system, user, assistant), and dispatches them to the local inference backend.
- **Streaming Response:** Text tokens stream to the client incrementally using standard Server-Sent Events (`text/event-stream`), enabling immediate reading without waiting for generation to finish.
- **Cancellation:** If the user cancels generation or closes the browser tab, generation terminates without blocking subsequent requests.
- **Backend Error Resilience:** If the local inference engine is unreachable or fails during generation, Linden emits structured, sanitized error states to the UI without logging user message content.

Originating PRD: [PRD-0001: Local Chat](../prd/0001-local-chat.md)

## Model selection

Linden queries the local inference daemon at startup and on demand (`GET /models`), allowing users to dynamically select from installed GGUF models directly in the web UI header.

## OpenAI API compatibility

External developer tools, IDE extensions, or agent scripts can target Linden's `POST /v1/chat/completions` endpoint as an alternative to cloud OpenAI endpoints, receiving streaming chat completions using standard OpenAI-formatted chunks.

## Zero-config LAN discovery

Linden broadcasts service presence across the local network via mDNS as `linden.local` (`_linden._tcp`), allowing any device on the same home network to connect without knowing the host's numerical IP address.
