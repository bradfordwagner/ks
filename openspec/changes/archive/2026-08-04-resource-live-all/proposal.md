## Why

`resource_live` already accepts an `-A/--all` flag to watch selections across all namespaces, but the interactive `ks` menu (no-argument invocation) only ever calls a command's `RunE` with default flags — it has no way to pass `-A`. `resource` solves this same problem for `resource_all` by existing as its own verb; `resource_live` has no equivalent, so all-namespaces live-scanning is unreachable from the menu.

## What Changes

- Add a new `resource_live_all` command that behaves exactly like `resource_live` but watches every selected type across all namespaces (equivalent to `resource_live -A`).
- Register `resource_live_all` in the interactive `ks` menu's command list so it is selectable without flags.
- Document the new verb in `CLAUDE.md`'s command reference table.

## Capabilities

### New Capabilities
(none)

### Modified Capabilities
- `resource-usage-filtering`: adds a requirement that `resource_live_all` performs the same live scan as `resource_live` but across all namespaces, mirroring how `resource_all` relates to `resource`.

## Impact

- `cmd/ks/resource_live_all.go` (new): registers the `resource_live_all` cobra command, calling `cmds.ResourceLive(standardArgs, true)`.
- `cmd/ks/main.go`: add `resourceLiveAllCmd` to the `commands` slice so it appears in the interactive menu.
- `CLAUDE.md`: add `resource_live_all` row to the Command Reference table.
- No changes to `internal/cmds`, `internal/resources`, or `internal/k9s` — `cmds.ResourceLive` already supports the `all` parameter used by `resource_live -A`.
