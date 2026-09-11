# PRD-0002: Local Document Ingestion & Grounded Retrieval (Local RAG)

| Field | Value |
|-------|-------|
| Status | Accepted |
| Author | Linden Architecture Team |
| Created | 2026-09-11 |
| Spec Reference | Planned for `docs/product/spec/` upon implementation |
| Supersedes | None |

---

## 1. Problem Statement & User Value

Users possess valuable, highly sensitive documents: tax returns, medical summaries, legal agreements, personal journals, and technical device manuals. 

Today, getting AI answers about personal documents requires either:
1. Uploading private files to cloud AI providers (violating privacy, risking data leaks, and granting third parties access to personal history).
2. Complex developer tooling that requires command-line expertise, Python environments, and external vector database subscriptions.

Linden solves this by bringing private, local Document Ingestion and Grounded Retrieval (RAG) directly into the home appliance. Users drop documents into Linden and get instant, grounded conversational answers with transparent citations—running privately on their local hardware.

---

## 2. Target Users & User Stories

- **Privacy-Conscious Homeowner:**
  - *Story:* As a homeowner, I want to index my mortgage and home appliance manuals so I can ask questions like "How do I clear error code E3 on my dishwasher?" without searching through filing cabinets or uploading manuals to external services.
- **Independent Professional / Researcher:**
  - *Story:* As an independent researcher, I want to query a repository of research notes and PDF papers offline so I can synthesize findings during travel or air-gapped field work.
- **Family / Household Member:**
  - *Story:* As a household member on the local network, I want clear, cited answers where clicking a source excerpt shows me exactly which document and page the AI used to answer my question.

---

## 3. Scope & Boundaries

### Must Have (Track C Delivery)
- **Local File Formats:** Direct ingestion and text extraction for Plain Text (`.txt`), Markdown (`.md`), and standard portable document format (`.pdf`).
- **Local Embeddings:** Automatic vector generation using dedicated local embedding models (e.g. `nomic-embed-text`) managed transparently by the local inference backend.
- **Grounded Chat Mode:** A toggle in the chat interface allowing users to query their document library or conduct general ungrounded chat.
- **Verifiable Citations:** Responses generated from retrieved context include visible source pill citations showing the document name and matched text excerpt.
- **Document Management:** Dedicated document library screen showing indexed files, file sizes, chunk counts, and indexing timestamps.
- **Permanent Local Purge:** Ability to delete any document, immediately removing its text and vector representations from disk.
- **System Concurrency & Resource Guardrails:** Background indexing throttles worker concurrency and memory usage so active chat streams remain responsive.

### Out of Scope (Non-Goals for Track C)
- **OCR for Scanned Images:** Extraction from scanned bitmap image PDFs without embedded text layers (deferred to future vision capabilities).
- **Proprietary Binary Formats:** Word (`.docx`), Excel (`.xlsx`), or PowerPoint (`.pptx`) parsing (plain text, markdown, and text-based PDFs first).
- **External Vector Clouds:** No connections to Pinecone, Weaviate, or remote embedding endpoints.
- **Cross-Network Synchronization:** Syncing document indices across multiple physical Linden devices.

---

## 4. High-Level Requirements

1. **R1 (Ingestion):** Users can add documents via drag-and-drop or file upload in the web interface.
2. **R2 (Parsing & Indexing):** Ingested files are automatically parsed, split into manageable semantic passages, and indexed into a persistent local vector store.
3. **R3 (Multi-Model Routing):** The platform transparently directs vector generation to an installed local embedding model without changing the user's selected conversational chat model.
4. **R4 (Grounded Querying):** When document search is active, user prompts retrieve relevant passages and synthesize an answer grounded strictly in the source materials.
5. **R5 (Source Citations):** Responses utilizing retrieved passages include citation links displaying the source file name and excerpted passage.
6. **R6 (Document Library & Inspection):** Users can browse all indexed documents, inspect their processing status (Ready, Indexing, Failed), and view metadata.
7. **R7 (Data Deletion):** Deleting a document purges the source file and deletes all corresponding index records permanently from the host.
8. **R8 (Resource & Queue Protection):** Heavy indexing jobs run with lower priority than interactive chat streams, ensuring low time-to-first-token during concurrent usage.
9. **R9 (Dependency Check):** If an embedding model is missing, Linden displays a clear, actionable guidance modal explaining which model to install.

---

## 5. Privacy, Safety, & Trust Constraints

- **Air-Gapped Operation:** All document parsing, embedding creation, vector similarity math, and prompt synthesis occur entirely within the host machine. Zero network bytes leave the LAN.
- **Strict Storage Boundaries:** Documents and vector embeddings are stored within the application's local user data directory with restrictive file permissions.
- **Zero Content Logging:** Document contents, file names, vector embeddings, and retrieval query snippets are strictly excluded from server operational logs.
- **Prompt Safety & Grounding:** System instructions instruct the model to state when an answer is not supported by the retrieved passages rather than inventing false facts.

---

## 6. Success & Acceptance Criteria

- **Functional Accuracy:** Asking a factual question directly answered in an ingested 30-page PDF returns the correct answer and cites the source document within 3 seconds on standard baseline hardware.
- **Zero Configuration:** A user can upload a markdown or PDF file and immediately begin querying it without manually configuring vector databases, chunk sizes, or API tokens.
- **Offline Integrity:** The entire lifecycle (upload → chunk → embed → query → cite → delete) functions identically with WAN network disconnects.
- **Memory Safety:** Ingesting a large document does not exceed process memory bounds or crash running services.

---

## 7. Open Questions

- *Storage Engine Selection:* Evaluate embedded SQLite with vector extensions (`sqlite-vec`) versus pure-Go in-memory/file vector index for simplicity and portability across Linux/macOS/Windows binaries.
- *Default Retrieval Top-K:* Benchmark response quality and prompt token consumption across top-3 vs. top-5 retrieved chunks for small local model context windows.
