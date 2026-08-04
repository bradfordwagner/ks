## Context

See proposal.md - Why. `resource_live` already supports all-namespaces scanning via its `-A/--all` flag (`cmds.ResourceLive(a, all bool)`); the gap is purely that the flag is unreachable from the interactive `ks` menu, which only invokes commands by name with default flags.

## Goals / Non-Goals

**Goals:**
- Make all-namespaces live-scanning reachable from the interactive menu.

**Non-Goals:**
- Any change to scan behavior, streaming, or prioritization logic in `cmds.ResourceLive` — it already takes an `all bool` and is reused unmodified.

## Decisions

- **Mirror `resource`/`resource_all` exactly**: `resource_live_all` is a new cobra command in its own file (`cmd/ks/resource_live_all.go`) whose `RunE` calls `cmds.ResourceLive(standardArgs, true)`, registered in `main.go`'s `commands` slice next to `resourceLiveCmd`. This is the established pattern in this codebase for exposing an all-namespaces variant to the flagless interactive menu — no alternative was considered, since deviating would introduce an inconsistent second pattern for the same problem `resource_all` already solves.

## Risks / Trade-offs

- [Duplication between `resource_live.go` and `resource_live_all.go`] → Both are ~15-line cobra wrappers with no logic of their own (all behavior lives in `cmds.ResourceLive`), matching the existing `resource.go`/`resource_all.go` duplication already accepted in this codebase.
