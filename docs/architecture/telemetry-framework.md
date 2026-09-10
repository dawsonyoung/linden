# Linden AI Telemetry & Debug Framework (Draft)

*Note: This is a conceptual draft to be revisited after the initial Linden MVP is completed.*

This document outlines the architecture for an "out-of-band" AI telemetry and diagnostic system for the Linden local AI platform.

## 1. Core Principles & Privacy Constraints

Given Linden is a **privacy-first, offline-first** platform, a telemetry system that communicates with developers introduces inherent friction. To solve this, the framework strictly adheres to the following:
- **Human-Readable Only:** The generated diagnostic payload is primarily plain English. It avoids proprietary or opaque binary data. Detailed technical logs are heavily sanitized.
- **Open Sanitization:** The code that sanitizes the telemetry payload is 100% open-source and publicly verifiable, ensuring trust in what is stripped out.
- **Explicit & Visible Transmission:** The system never sends data in the background. The user must review the payload and explicitly approve transmission. Furthermore, a permanent record of the transmission (e.g., an internal Linden log or an email confirmation) is generated for the user's peace of mind.
- **Offline Fallback:** If the device is air-gapped or offline, the human-readable report is presented to the user to copy and paste into an email or portal when they reach a connected device.

## 2. Architecture: The "Side Channel"

To ensure the diagnostic system survives failures in the main application flow, it requires an isolated execution environment.

**Recommendation:** A dedicated Diagnostic Subsystem running within the same binary, but on an isolated listener.

- **Why?** Shipping a single binary maintains Linden's "plug-and-play" simplicity.
- **How it works:** When `cmd/linden` starts, it spawns the main `api` server on the primary port (e.g., `:8080`). It also spawns the `telemetry` server on a restricted local Unix Domain Socket or a hidden localhost-only port (e.g., `127.0.0.1:8081`).
- **Isolation:** The telemetry system maintains its own instance of the `inference` client, completely bypassing the `orchestrator`. If the main orchestrator deadlocks or the API crashes, the side-channel remains fully responsive.

```mermaid
graph TD
    subgraph Single Linden Binary
        subgraph Main Application
            API[api/:8080] --> Orch[orchestrator/]
            Orch --> InfMain[inference/]
            Orch --> Store[storage/]
        end

        subgraph Out-of-Band Diagnostic Channel
            DiagAPI[telemetry/:8081 or UDS] --> Agent[Diagnostic AI Agent]
            Agent --> InfDiag[inference/ (Isolated Instance)]
            
            %% Observability hook
            Agent -.->|Read/Toggle State| Main Application
        end
    end
    
    User[User / Web UI] -->|Normal Traffic| API
    User -->|Report Issue| DiagAPI
```

## 3. Workflow & Data Flow (Client Side)

1. **Trigger:** The user encounters an issue and clicks "Diagnose Issue" in the UI.
2. **Investigation:** The Diagnostic Agent uses its isolated inference client to investigate. It toggles internal observability hooks, reads recent error states, and evaluates the environment (e.g., Ollama status).
3. **Report Generation:** The Agent drafts a human-readable English report describing the issue, supplemented with sanitized snippets (stack traces). 
4. **User Review:** The user is presented with the complete, plaintext payload.
5. **Transmission:** 
   - **Online:** The user clicks "Send". The payload is sent to Linden HQ. A visible receipt is logged in the user's Linden UI.
   - **Offline:** The user copies the text to manually email it from a connected device.

## 4. The Automated Fix Pipeline (HQ Side)

When the report is received by Linden HQ (via email or direct API), an automated CI pipeline attempts to resolve the issue without human intervention.

1. **Ingestion & Triage:** An AI agent at HQ receives the payload. Because the report is in plain English with localized stack traces, the agent parses the context to identify the failing component (e.g., `orchestrator`, `storage`, `web`).
2. **Replication (Test-Driven):** Before altering any product code, the HQ agent attempts to write a failing unit or contract test that successfully replicates the user's reported bug.
3. **Automated Implementation:** The agent modifies the source code to fix the bug, ensuring the new test passes.
4. **Validation:** The full suite of quality gates (`make validate`) is executed, running all contract, integration, and security tests to ensure the fix didn't introduce regressions.
5. **Review & Release:** 
   - A pull request is generated automatically. 
   - Depending on confidence scores, the PR is either merged automatically or flagged for human review.
   - Once merged, a specific patch build (e.g., `v1.2.4-patch.1`) is compiled.
6. **User Notification:** The user receives a notification (either via email or their local Linden UI) stating: *"A potential fix for your reported issue is ready. Would you like to download test build v1.2.4-patch.1 to verify it resolves your problem?"*

## 5. Structural Additions (`src/`)

To support this cleanly without violating Linden's dependency rules (when implemented):
- **`src/telemetry/`**: The new layer responsible for the diagnostic agent and sanitization logic. Depends on `inference/` and `errs/`.
- **`src/api/` (Observability Interface)**: Exposes internal state variables (expvar, pprof, log levels) that `src/telemetry/` can read.
- **`src/cmd/main.go`**: Wires up both the main `api` server and the isolated `telemetry` listener.
