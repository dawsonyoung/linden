# Linden Long-Term Strategic Roadmap

*Note: This document compiles the long-term architectural, commercial, and operational roadmap (Tracks D through H) for external business planning and strategic reference.*

---

## 1. Track D: Security, Privacy, and Trust Controls

### D1. Security Baseline
1. **Local TLS (HTTPS)**: Automatic local certificate provisioning or self-signed CA generation for secure LAN transmission.
2. **Session Hardening**: Cryptographic session tokens for client pairing, automatic token expiration, and rotation.
3. **Local Secrets Vault**: Secure on-disk storage for optional API keys or tokens using OS keychain or local encryption.
4. **Hardened Defaults**: Strict defaults with clear UI warnings for high-risk network configurations.

### D2. Privacy Controls
1. **Data Retention Policies**: User-selectable presets (e.g. ephemeral session-only, local persistent history, custom).
2. **Per-Feature Outbound Toggles**: Granular network gates explaining exactly what external services are connected (if any).
3. **Local Audit Log**: Searchable in-browser viewer showing all model invocations and tool executions without logging sensitive text.
4. **Data Portability**: One-click JSON/zip export of all local conversations and indexed documents, with permanent local purge capability.

### D3. Verification & Assurance
1. **Threat Model Updates**: Documented attack surface analysis covering local prompt injection and LAN traversal.
2. **Automated Security Suites**: Automated regression packs in `validation/security`.
3. **Provenance Reporting**: Signed release binaries and verifiable container digests.

---

## 2. Track E: Product Experience and Multi-Device UX

### E1. First-Run Experience
1. **Setup Wizard**: Automated hardware detection and model download recommendation flow.
2. **Network Onboarding**: Plain-English explanation of local vs. remote operation during initial launch.
3. **Recovery Flows**: Automated diagnosis and guided resolution for missing dependencies (e.g. Ollama offline).

### E2. Daily Use UX
1. **Conversation Management**: Chat organization with search, pinning, and tagging.
2. **Prompt Templates**: Reusable local system prompts and custom assistant personas.
3. **Explainability Panel**: Visual inspection pane showing sources cited and tool actions executed during a reply.

### E3. Multi-Device LAN Access
1. **Frictionless Pairing**: Simple LAN approval handshake for secondary devices (phones, tablets, laptops).
2. **Device Roles**: Access controls distinguishing administrator devices from client devices.
3. **Revocation & Anomaly Alerts**: Instant pairing revocation for untrusted or obsolete devices.

---

## 3. Track F: Developer Ecosystem and Extensibility

### F1. SDK & Interface Contracts
1. **Client Libraries**: Official typed clients for web, mobile, and CLI environments.
2. **Extension Policy**: Versioned contracts ensuring backward compatibility for plugins and tools.
3. **Developer Sandbox**: Reference implementations and testing harness for custom tool development.

### F2. Tooling Ecosystem (Model Context Protocol - MCP)
1. **MCP Host Support**: Built-in support for hosting and connecting local MCP servers over stdio and HTTP/SSE.
2. **First-Party Tools**: Curated tools for local filesystem search, calendar, notes, and Home Assistant IoT.
3. **Community Tool Signing**: Signature verification ensuring third-party extensions cannot perform unverified network calls.

### F3. Workflow Automation
1. **Trigger-Based Workflows**: Scheduled or event-driven local agent automations.
2. **Approval Checkpoints**: Mandatory human-in-the-loop gates for any tool call with side effects (writing files, making network requests).

---

## 4. Track G: Hardware Appliance Strategy

### G1. Appliance MVP
1. **Target Hardware SKUs**: Low-power x86-64 mini-PCs (e.g. Intel N100 / AMD Ryzen) and high-memory ARM SBCs.
2. **First-Boot Appliance Image**: Minimal Linux OS image configured to boot directly into Linden on LAN connect.
3. **Diagnostics Partition**: Dedicated recovery environment for offline appliance health checks and factory resets.

### G2. Manufacturing & Fleet Operations
1. **Secure Boot & Firmware**: Signed boot chain protecting appliance integrity.
2. **Appliance Telemetry Boundaries**: Strict local-only hardware health monitoring with opt-in support bundles.
3. **Appliance-to-Software Migration**: Simple migration utility allowing users to transfer appliance state to any PC.

---

## 5. Track H: Business, Licensing, and Governance

### H1. Licensing Operations
1. **Source-Available Boundary**: Clear distinction between personal/evaluation use and commercial/enterprise hosting.
2. **Contributor Agreement**: Clear terms for open-source and source-available contributions.

### H2. Commercialization Model
1. **Paid Convenience Features**: Monetization through turnkey hardware appliances, specialized pre-indexed knowledge packs, or optional managed remote tunnels—never through user data monetization.
2. **Team & Small Business Packs**: Local-first multi-user permissioning and governance for privacy-sensitive workplaces.

### H3. Trust Governance
1. **Public Privacy Changelogs**: Mandatory public disclosure for any change affecting network boundaries or storage.
2. **Security Vulnerability SLA**: Published disclosure timelines and rapid response processes.
