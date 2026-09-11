# Model Context Protocol (MCP) Boundary Specification

## 1. Overview & Boundary

The `mcp` layer (`src/mcp/`) establishes Linden's architectural boundary for extensibility via the Model Context Protocol (MCP). It defines how local tools, custom integrations, and data resources are discovered, sandboxed, and invoked by the orchestrator.

- **Package:** `github.com/dawsonyoung/linden/mcp`
- **Dependencies:** Leaf layer boundary. Depends only on `errs`. No direct dependencies on `api` or `cmd`.
- **Primary Consumer:** Consumed by `src/orchestrator/` during agent tool-dispatch and retrieval workflows.

---

## 2. MCP Wire Protocol & Message Types

Linden adheres to the JSON-RPC 2.0 transport specification for MCP tool communication over `stdio` or local UNIX domain sockets:

| Method | Direction | Purpose | Schema |
|---|---|---|---|
| `initialize` | Client ──► Server | Negotiates capabilities and protocol versions. | Protocol version handshake. |
| `tools/list` | Client ──► Server | Discovers tools exposed by the MCP server. | Returns list of `ToolDefinition`. |
| `tools/call` | Client ──► Server | Executes a specific tool with arguments. | Accepts `name` and `arguments` object; returns `content` blocks. |
| `resources/list` | Client ──► Server | Discovers readable local resources. | Returns resource URIs and MIME types. |
| `resources/read` | Client ──► Server | Fetches data from a resource URI. | Returns text or binary content blob. |

---

## 3. Tool Manifest Schema Table

| Field | Type | Required | Constraints | Description |
|---|---|---|---|---|
| `name` | `string` | Yes | 1–64 alphanumeric chars | Unique tool identifier (e.g. `"search_local_docs"`). |
| `description` | `string` | Yes | 10–500 chars | Explains tool behavior to the LLM. |
| `inputSchema` | `object` | Yes | Valid JSON Schema | Formal schema validating tool input parameters. |
| `outputSchema`| `object` | No | Valid JSON Schema | Schema validating tool output results. |

---

## 4. Execution Sandbox & Security Guardrails

To protect user privacy and device stability on the home network, all MCP tool executions must satisfy strict sandboxing constraints:

| Guardrail | Constraint | Violation Action |
|---|---|---|
| **Process Isolation** | Tool plugins execute in isolated child processes with restricted user permissions. | Subprocess abort. |
| **Timeout Budget** | Default 15-second execution ceiling per tool invocation. | `errs.DeadlineExceeded` and SIGKILL. |
| **Output Size Cap** | Maximum 512 KB returned text payload per tool call. | Output truncated; error returned. |
| **Network Egress** | Zero ambient outbound internet network access without explicit user consent. | TCP socket blocked. |
| **File Sandbox** | File access restricted exclusively to the designated local workspace. | Path traversal rejected with `errs.PermissionDenied`. |

---

## 5. Error Taxonomy & Status Mapping

| MCP Error Condition | Mapped `errs.Code` | Notes |
|---|---|---|
| Unknown tool name in `tools/call` | `errs.NotFound` | Tool not registered in manifest. |
| Arguments violate `inputSchema` | `errs.InvalidArgument` | Schema validation failed before invocation. |
| Tool execution exceeded timeout | `errs.DeadlineExceeded` | Process killed after timeout. |
| Tool process crashed / panic | `errs.Internal` | Subprocess exited with non-zero code. |
| Tool attempted unauthorized access | `errs.PermissionDenied` | Sandbox boundary violation. |
