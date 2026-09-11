# MCP (Model Context Protocol) Layer

This package is reserved for Linden's Model Context Protocol (MCP) integrations.

## Future Scope

As Linden evolves from a standalone chat appliance into a robust, extensible local AI platform, this directory will house:
- **MCP Client Integrations**: Allowing Linden's underlying orchestrator to seamlessly query local tools (e.g., local filesystem readers, local web search, Home Assistant IoT).
- **MCP Server Capabilities**: Enabling Linden to expose local models to other MCP-compliant applications.

**Note:** Code in this directory should strictly adhere to Linden's boundary rules. It should expose clean Go interfaces to the `orchestrator` layer without bleeding protocol-specific JSON logic outside of `src/mcp/`.
