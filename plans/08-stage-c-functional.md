# Stage C: Functional and Hardening

## Status
Ready to implement.

## Purpose

Implement the user-facing frontend, deploy the full end-to-end stack in a Docker container, and secure the system against local vulnerabilities and data leakage. This stage transitions Linden from a functioning API to a usable, secure product.

Branching rules: `docs/project/branching-strategy.md`.
Delivery policy: `docs/project/contract-first-delivery.md`.

## Sequencing

```
C.1 feat/web-chat-ui ──→ C.2 test/smoke-container ──→ C.3 test/security-baseline
```

---

## C.1 `feat/web-chat-ui`

Intent: Create the user-facing chat application.

Agents: Web Agent, Implementation Agent, Review Agent.

### Tasks

1. Initialize a SvelteKit project in `src/web/` using the static adapter.
2. Install `tailwindcss`, `daisyui`, and `@ai-sdk/svelte` for the UI styling and chat streaming implementation.
3. Build the chat interface capable of capturing user input, displaying historical turns, and rendering streamed markdown responses from the assistant.
4. Ensure the static output of the web build is served from the `src/api` layer (via Go's `embed` or serving the build directory).

### Gate

Lint, Web tests (SvelteKit checks).

Estimated size: 300–500 lines.

---

## C.2 `test/smoke-container`

Intent: Validate the end-to-end functionality of the fully containerized application.

Agents: Docker Agent, Integration Agent, Review Agent.

### Tasks

1. Update the `api` listener to bind to `0.0.0.0` to permit Local Area Network (LAN) testing from separate devices.
2. Ensure `docker-compose.yml` mounts the correct storage volumes for sessions and binds the correct ports.
3. Add a smoke test script that starts the container, waits for `/health`, and verifies the web UI is served correctly.
4. Provide instructions on testing the UI from a separate device on the local network.

### Gate

Docker build, Smoke tests pass.

Estimated size: 100–150 lines.

---

## C.3 `test/security-baseline`

Intent: Ensure the privacy-first model is strictly enforced before MVP hardening.

Agents: Security Agent, Contract Test Agent, Review Agent.

### Tasks

1. Write tests that simulate malformed inputs (e.g., prompt injections, massive payloads) to ensure the API safely rejects them.
2. Review and ensure no user chat data or sensitive information is printed to stdout or logs.
3. Ensure error responses returned to the client are generic and sanitized (no stack traces).

### Gate

Security tests pass.

Estimated size: 150–200 lines.

## Stage C Exit Criteria

1. The user can navigate to the local IP address on any device and access a polished SvelteKit UI.
2. The UI can send messages and stream responses flawlessly via `@ai-sdk/svelte`.
3. The Docker container provisions successfully, exposing the API and UI safely.
4. No sensitive user chat content is leaked to logs or error messages.
