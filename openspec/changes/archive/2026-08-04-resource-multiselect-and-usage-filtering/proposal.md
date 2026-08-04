## Why

Watching a spread of an operator's CRDs (e.g. the ~20 vault-config-operator types) one at a time is tedious, and once the catalog is group-qualified the remaining noise is CRD types that have no objects. This change makes the resource picker multi-select (tiling a k9s pane per pick) and adds two ways to narrow the picker to resource types that actually hold objects.

This supersedes the earlier `resource-multi-watch` change: rather than a separate `resource_multi` command, multi-selection is folded into `resource`/`resource_all`.

## What Changes

- `resource` / `resource_all` become **multi-select**: one pick opens k9s in the current pane (unchanged); multiple picks open the first in the current pane and one tiled split pane per remaining selection, all sharing the current `KUBECONFIG`.
- **BREAKING (unshipped):** the `resource_multi` command is removed — its behavior is folded into `resource`.
- `resource_load` also **probes usage**: a parallel `Limit:1` existence check per listable type (all namespaces), recording `populated`/`probed` per type in `resources.json` (schema v3).
- New `resource_exists`: multi-select from types recorded as populated at the last `resource_load` — instant, no live calls.
- New `resource_live`: live-scan the cluster and multi-select only types that currently have objects, ordered so previously-populated + voted types come first.

## Capabilities

### New Capabilities
- `resource-selection`: how the resource picker selects and launches — multi-selection that opens the first pick in the current pane and a tiled split pane per additional pick.
- `resource-usage-filtering`: probing resource types for objects at load time and filtering the picker to populated types, cached (`resource_exists`) or live (`resource_live`).

### Modified Capabilities
<!-- None. KS_RESOURCE/cache fast-path behavior (ks-resource-env-override) and group-qualified naming (group-qualified-crd-names) are unchanged and referenced, not modified. -->

## Impact

- **Code**: `internal/cmds/resource_cmd.go` (multi-select + `launchSelected`/`runSingle`/`singleResource` helpers), `internal/cmds/discovery.go` (`listableGVRs`, `probePopulated`), `resource_load_cmd.go` (probe + flags), new `resource_live_cmd.go` / `resource_exists_cmd.go`, `internal/tmux/split.go` (`SplitAndSend`), `internal/resources/types.go` (`Populated`/`Probed`, schema v3).
- **Commands**: `resource_multi` removed; `resource_live`, `resource_exists` added in `cmd/ks` + `main.go` menu.
- **Data**: `resources.json` gains `populated`/`probed` per entry (additive; v1/v2 load fine).
- **Runtime**: adds use of the dynamic client for probing; tmux-only for the split panes (same as before).
