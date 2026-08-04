## Why

None of `ks`'s fzf pickers (`internal/choose`) indicate which cluster/context is currently active, making it easy to lose track of what's selected when several terminals are open against different clusters. This was identified during a feature-gap review as a small, isolated UX fix.

Note: an unfiltered-kubeconfig-listing fix (`internal/list.Kubeconfigs` picking up junk files like `.DS_Store`/`*.bak`) was originally scoped into this same change and implemented, but was walked back at the user's request to keep this change focused on the context indicator only. It may be proposed separately later.

## What Changes

- `internal/choose.One` and `internal/choose.Multi` gain a required `header` parameter that is rendered as a `[<context>] > ` prompt prefix (via `fzf.WithPrompt`), colored Catppuccin Mocha blue (`#89b4fa`, via `fzf.WithStylePrompt`), instead of the bare uncolored `> ` default. **BREAKING** (internal API only, not user-facing): all existing callers of `choose.One`/`choose.Multi` must be updated to pass a header argument.
- New `internal/kube.CurrentContext(kubeconfig string) string` resolves the current-context name from a kubeconfig file, falling back to the base filename if it can't be read.
- Every command that shows an fzf picker (root menu, `set_ns`, `resource`/`resource_all`, `kube_new_ns`, `tmux_multi`, `tmux_window`, `kube_cp`, `link`, `pipe`) passes `kube.CurrentContext(a.Kubeconfig)` as the header so the active cluster/context is always visible above the prompt.

## Capabilities

### New Capabilities
- `picker-context-indicator`: every interactive fzf picker in `ks` displays the currently active kubeconfig context as part of its prompt.

### Modified Capabilities
(none — no existing spec covers picker behavior)

## Impact

- Affected code: `internal/choose/choose.go`, `internal/kube/client.go`, `cmd/ks/main.go`, and every file under `internal/cmds/` that calls `choose.One`/`choose.Multi` (`set_ns_cmd.go`, `resource_cmd.go`, `kube_new_ns_cmd.go`, `tmux_multi_cmd.go`, `tmux_window_cmd.go`, `kube_cp_cmd.go`, `link_cmd.go`, `pipe_cmd.go`).
- Tests: `internal/cmds/resource_cmd_test.go` updated for the new `resolveResourceType` signature.
- No external dependency changes; `go-fzf@v0.15.0`'s existing `WithPrompt` option is reused.
