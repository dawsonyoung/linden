# Security Policy

## Reporting a vulnerability

Report privately. Do not open a public issue.

Use [GitHub private vulnerability reporting](https://github.com/dawsonyoung/linden/security/advisories/new) for this repository.

Include what you can: affected version or commit, reproduction steps, and impact. A partial report is more useful than no report.

Expect acknowledgement within 7 days. Linden is currently maintained by one person, so remediation time depends on severity and scope. You will be told which applies.

## Scope

Linden is pre-release. No version is supported for production use, and there is no security-fix backport channel yet.

| Version | Supported |
|---------|-----------|
| `main` | Best effort |
| Tagged releases | None yet |

## Security model

Linden is designed to run on hardware the user controls, on a trusted local network. Assumptions:

- The host machine is trusted.
- The local network is semi-trusted; LAN access requires explicit pairing.
- The model backend runs locally.
- No outbound connection is made unless the user enables one.

Findings that break these assumptions are in scope. A finding that requires an already-compromised host generally is not.

Published privacy commitments are in [docs/product/spec/30-privacy-model.md](docs/product/spec/30-privacy-model.md). A deviation between that document and actual behavior is a security issue, not a documentation bug — report it here.

## Standing rules

These are enforced in review and are treated as defects when violated:

- No user content in logs, at any level, including debug.
- No user content in error messages returned to clients.
- All input validated at the API boundary.
- No hardcoded secrets, tokens, or keys.
- No user data leaves the device without explicit user action.

## Dependencies

Linden targets the Go standard library. Every external dependency requires justification in the PR that introduces it. Fewer dependencies means less to audit.
