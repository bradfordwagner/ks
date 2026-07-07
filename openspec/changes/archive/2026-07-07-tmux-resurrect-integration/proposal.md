## Why

When tmux-resurrect restores a session, pane IDs and environment variables are not preserved — so `KUBECONFIG` is lost and the per-pane resource cache in `.ks.resources.json` becomes stale. Users must manually re-run `ks tmux_window`/`ks tmux_multi` and re-select their resource for every pane after every resurrect.

## What Changes

- New `ks save` command: scans all tmux panes at resurrect-save time, extracts `KUBECONFIG` and cached resource per pane via `/proc`, writes a positional sidecar file `~/.kube/.ks.resurrect.json`
- New `ks restore` command: reads the sidecar, matches panes by stable positional index, sends `KUBECONFIG` and re-launches `ks resource` / `ks resource_all` into each pane
- New `KS_RESOURCE` env var support in `Resource()`: allows bypassing fzf selection for a single invocation without persisting into the shell environment
- New sidecar file type: `~/.kube/.ks.resurrect.json` keyed by `session/window_idx/pane_idx`

## Capabilities

### New Capabilities

- `resurrect-save`: Scan all tmux panes, extract KUBECONFIG and resource from /proc and the cache, write positional sidecar for resurrect
- `resurrect-restore`: Read sidecar, match panes by position, re-apply KUBECONFIG and re-launch ks resource/resource_all inline
- `ks-resource-env-override`: `KS_RESOURCE` env var bypass in `Resource()` — single-invocation, not exported, preserves clear_pane/clear_cache semantics

### Modified Capabilities

## Impact

- `internal/cmds/resource_cmd.go`: add `KS_RESOURCE` env var check before cache lookup
- `internal/cmds/` (new): `resurrect_save_cmd.go`, `resurrect_restore_cmd.go`
- `internal/resurrect/` (new): sidecar types, proc reading, tmux pane listing
- `cmd/ks/main.go`: register `save` and `restore` Cobra commands
- New file: `~/.kube/.ks.resurrect.json` (runtime, not committed)
- Linux/WSL2 only for save (`/proc` dependency); restore works anywhere
