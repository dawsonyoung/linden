# PRD-0001: Local Streaming Chat

| Field | Value |
|-------|-------|
| Status | **Shipped** |
| Author | dawsonyoung |
| Created | 2026-08-15 |
| Spec Reference | [docs/product/spec/10-capabilities.md](../spec/10-capabilities.md) |
| Supersedes | None |

---

## 1. Problem Statement & User Value

Users who want AI assistance today are forced to send their personal thoughts, questions, and sensitive data to remote cloud services. There is no way to audit what is retained or reused for AI training, and assistants completely stop functioning when an internet connection drops. 

Linden solves this by providing a personal conversational assistant that runs 100% on the user's home network, keeping every byte of prompt and response text on local hardware.

## 2. Target Users & User Stories

- **Privacy-Conscious Homeowner:** As a home user, I want to ask questions and draft sensitive text without my input ever leaving my private premises.
- **Offline Operator:** As a user on an air-gapped or unreliable network, I want a responsive assistant that never breaks during internet outages.

## 3. Scope & Boundaries

### Must Have
- Send conversational messages and receive replies generated entirely by a local LLM backend.
- Incremental streaming text delivery over standard Server-Sent Events (SSE) so users don't wait for complete generation.
- User cancellation: ability to stop generation mid-stream without crashing or blocking subsequent prompts.
- Dynamic local model inspection and selection via API endpoints.
- Actionable error feedback when the local model engine is offline or missing.

### Out of Scope (Initial Release)
- Personal document ingestion and vector search (RAG) — deferred to Track C / PRD-0002.
- Multi-device pairing handshakes and access roles — deferred to Track E.
- Remote cloud model fallback providers.

## 4. High-Level Requirements

1. **R1:** Messages sent from the client receive replies generated entirely on local hardware.
2. **R2:** Reply text appears incrementally in the interface as tokens are generated.
3. **R3:** The user can stop generation at any time, retaining the partial reply received so far.
4. **R4:** If no local model is available or installed, Linden displays an actionable notice explaining how to pull a model.
5. **R5:** If the local inference engine is unreachable, the client remains responsive and presents an actionable recovery message.
6. **R6:** The system can route requests to any locally installed model via API parameters (interactive web UI selector planned for subsequent UI cycle).

## 5. Privacy, Safety, & Trust Constraints

- **In-Memory Transport:** Message tokens are dispatched to the local inference daemon over loopback and held in memory for request duration.
- **Strict Zero-Logging:** Prompt text and generated tokens are strictly forbidden from server logs at all levels (including `DEBUG`).
- **Local Network Privacy:** During local operation, no tokens, user telemetry, or identifiers are transmitted outside the local home network.

## 6. Success & Acceptance Criteria

- A user can plug in the device, navigate to `http://linden.local:8080`, and begin a streaming chat session with zero configuration.
- Disconnecting the WAN internet connection does not alter chat functionality or responsiveness.
