# Linden System Architecture

This document outlines the core software architecture for the Linden local AI appliance. 

## 1. Software & Runtime Stack

### Inference Layer
- **Engine:** Headless `llama.cpp` server or `Ollama` daemon running on CPU with AVX2/VNNI acceleration.
- **Model Targets (Quantized GGUF):** 
  - Primary: Qwen 2.5 7B (Q4_K_M) or Llama 3.2 3B.
  - Embeddings: `nomic-embed-text` or `bge-m3` running locally.
- **Tuning Targets:** CPU thread pinning (`-t 4` matching physical cores), context window sizing capped to protect the 16GB RAM ceiling.

### Core Orchestration Layer (Custom Go Engine)
- **Language:** Go (single static binary, scratch/Alpine container, <30MB RAM footprint).
- **Protocols & Streaming:**
  - Exposes an OpenAI-compatible Server-Sent Events (SSE) endpoint (`/v1/chat/completions`) matching Vercel AI SDK conventions.
  - Programmatic ReAct loops, deterministic state machines, and typed schema validation for tool calls.
- **Local RAG & Storage:**
  - Document parsing (text, PDF, markdown) for personal finance, medical records, and education materials.
  - Local vector retrieval using embedded `sqlite-vec` or hybrid BM25 + dense search on-disk.

### Tooling & Extensibility (Model Context Protocol - MCP)
- **Architecture:** Tools written as standalone Go MCP servers communicating over `stdio` and HTTP/SSE transports.
- **Portability:** Decoupled from the UI; capable of being queried by any MCP client or the internal Go agent loop. This forms the technical foundation for the ecosystem marketplace.

### Content Safety & Filtering (The "Den" Persona)
- **Input/Output Guardrails:** Fast Go regex/Aho-Corasick PII scrubbing; optional local micro-classifier (Meta Prompt-Guard / Llama Guard 1B/3B) as middleware before tool dispatch.
- **Execution Rails:** Strict filesystem path traversal checks, schema whitelisting, and read-only boundaries.

---

## 2. Frontend & User Experience
- **Stack:** Lightweight Next.js / React client utilizing `@ai-sdk/react` (`useChat`) and `assistant-ui` / shadcn AI primitives.
- **Deployment:** Served directly over local network (`http://linden.local` via mDNS/Avahi) for a zero-configuration "plug-and-play" experience.
- **Key Features:** Streaming code blocks, markdown parsing, tool call expanders, drag-and-drop document upload for local RAG.

---

## 3. OTA Update Pipeline (Alpha to Beta)
- **Decoupled Updates:**
  - *Binary/UI:* Signed ~20–30MB Go binaries delivered via Cloudflare R2 + Worker manifest API.
  - *Models:* Asynchronous chunked GGUF downloads.
- **Safety Mechanism:** Atomic symlink binary replacement (`/opt/app/bin/current`), Ed25519 signature verification, and 60-second systemd watchdog self-test with automatic rollback to previous version on panic/crash to prevent appliance bricking.
