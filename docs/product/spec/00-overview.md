# Overview

## What Linden is

Linden is a privacy-first local AI platform. A user runs it on a machine they own, on their own network, and gets a personal AI assistant. There are no accounts, no cloud dependency, and no telemetry.

## Design commitments

These are guarantees, not aspirations. Each one is testable and each one constrains the implementation.

1. **Local by default.** No user content leaves the device unless the user explicitly enables an outbound path.
2. **Observable.** The user can see what data exists and remove it.
3. **Offline-capable.** Full functionality without internet access.
4. **Portable.** Data can be exported in a documented format.
5. **No hidden retention.** Nothing is stored that the specification does not describe.

## Scope of this document

This specification describes observable behavior — what a user or client can see and rely on. It does not describe internal structure; that is in `docs/architecture/layer-interface-spec.md`.

## Current state

Linden currently runs as a standalone, single-binary application providing:
- **Embedded Web Client:** A responsive SvelteKit interface compiled into the Go binary and served over HTTP.
- **Zero-Config LAN Discovery:** Advertises `linden.local` across the local Wi-Fi/LAN via mDNS.
- **Streaming Local Chat:** Multi-turn conversational chat with live token streaming over Server-Sent Events (SSE), backed by a local Ollama inference daemon.
- **Model Selection:** Dynamic inspection and selection of locally installed models via `GET /models`.
- **OpenAI Compatibility:** A standard `POST /v1/chat/completions` translation endpoint for drop-in integration with external tools and agent workflows.
- **Operational Endpoints:** Liveness checks (`GET /health`), build identification (`GET /version`), and request correlation headers (`X-Request-ID`).
