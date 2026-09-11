# Track C: Platform Scaling and Retrieval Maturity

## Purpose

Track C establishes Linden's data persistence and retrieval maturity (local RAG), platform hardening, concurrency handling, and multi-model capability routing. This work builds directly on top of the completed MVP Baseline (Track A) and Release Packaging (Track B).

## Scope & Milestones

### C1. Foundations & Platform Hardening
1. **API Namespace Evolution**: Standardize endpoints under versioned routing while maintaining backwards-compatible root endpoints.
2. **Configuration Profiles**: Add profile presets (`basic`, `power-user`, `hardened`) to control memory ceilings and resource consumption.
3. **Data Schema Migrations**: Implement a lightweight migration runner for local SQLite/JSON storage schemas to protect user data across binary upgrades.

### C2. Reliability and Performance
1. **Soak and Concurrency Testing**: Validate that concurrent requests from multiple LAN clients queue safely without deadlock or memory exhaustion.
2. **Latency & Buffer Optimization**: Streamline chunk serialization between the inference engine and the HTTP SSE writer.
3. **Request Budgeting & Backpressure**: Enforce graceful rejection or queueing when local CPU/RAM utilization exceeds safety thresholds.
4. **Memory Guardrails**: Establish proactive thread and context window caps based on host memory tier (<8GB, 16GB, 32GB).

### C3. Data and Retrieval Maturity (Local RAG)
1. **Document Ingestion Pipeline**:
   - Parsing adapters for local files: plain text, markdown, and PDF.
   - Text chunking strategies with source metadata preservation.
2. **Local Vector & Keyword Search**:
   - In-process or embedded vector retrieval (e.g. `sqlite-vec` or local embeddings via Ollama `nomic-embed-text`).
   - Hybrid lexical (BM25) and semantic search ranking.
3. **Context Assembly & Citations**:
   - Dynamic prompt context formatting with verifiable file name and chunk citations in chat stream responses.
   - User-configurable retrieval threshold and source toggle per session.

### C4. Multi-Model Capability
1. **Model Capability Registry**: Inspect local engine models and tag them by capability (`chat`, `embedding`, `code`, `vision`).
2. **Task-to-Model Routing**: Automatically route embedding jobs to dedicated embedding models while keeping chat focused on the user-selected conversational model.
3. **Fallback Chains**: Graceful degradation when high-parameter models fail or experience out-of-memory errors on constrained hardware.
