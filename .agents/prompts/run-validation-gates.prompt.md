---
name: "Run Validation Gates"
description: "Execute and interpret all quality gates. Use before packaging a PR or when a gate failure needs diagnosis."
---
# Run Validation Gates

## Gates and Commands

| Gate | Command | Pass Condition |
|------|---------|---------------|
| Lint | `cd src && go vet ./...` | Exit 0 |
| Unit | `cd src && go test -race ./...` | Exit 0 |
| Contract | `cd src && go test -tags=contracts -race ../validation/contracts/...` | Exit 0 |
| Integration | `cd src && go test -tags=integration -race ../validation/integration/...` | Exit 0 |
| Docker build | `docker build -t linden:ci .` | Exit 0 |
| Health check | `curl -sf http://localhost:8080/health` after container start | HTTP 200 |
| Web check | `cd src/web && npm run check` | Exit 0 |

## Procedure

1. Run each gate in order.
2. Record PASS or FAIL and the last 40 lines of output for any failure.
3. Stop at first BLOCKER failure and hand off to Debug Agent with the failure output.
4. On full pass, produce a gate evidence summary for the PR body.

## Output

Gate evidence table (gate name, status, notes) suitable for pasting into a PR description.
