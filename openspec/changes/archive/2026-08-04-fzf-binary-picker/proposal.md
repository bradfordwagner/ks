## Why

The picker (`internal/choose`) used the `go-fzf` Go library, which does not read the user's `FZF_DEFAULT_OPTS` or key binds, so ks pickers behaved differently from every other fzf in the user's shell. Shelling out to the real `fzf` binary fixes that: ks pickers inherit the user's fzf configuration and behave like every other fzf they use.

## What Changes

- `internal/choose` shells out to the **`fzf` binary** instead of `go-fzf`. Every picker (resource, kubeconfig, namespace, command menu) now inherits the user's `FZF_DEFAULT_OPTS` and key binds.
- Inside tmux, pickers render in a centered **`--tmux` popup** so tmux owns the popup's screen region — which fixes the full-screen alt-screen redraw corruption some terminals show (notably for incremental/streaming updates). Outside tmux, fzf renders normally.
- `One()` forces single-select (`--no-multi`) so it is correct even when the user's defaults enable multi (`-m`).
- New `MultiStream()` feeds fzf's stdin from a channel so options appear **live** as a producer emits them (used by `resource_live`); rendering cleanly relies on the popup.
- Picker abort semantics are centralized as `choose.ErrAbort` (returned on Esc/Ctrl-C = exit 130, and no-match = exit 1); real fzf errors are surfaced.
- The `github.com/koki-develop/go-fzf` dependency is removed.

## Capabilities

### New Capabilities
- `interactive-picker`: how ks presents interactive selection — backed by the `fzf` binary, honoring the user's fzf configuration, with single/multi/streaming modes and consistent abort handling.

### Modified Capabilities
<!-- None at the requirement level. The KS_RESOURCE fast-path (ks-resource-env-override) and resource selection/usage capabilities consume this picker unchanged. -->

## Impact

- **Code**: `internal/choose/choose.go` rewritten (`One`/`Multi`/`MultiStream`/`ErrAbort`, `--tmux` popup in tmux); all 11 callers (`cmd/ks/main.go` + 10 `internal/cmds/*`) switch `fzf.ErrAbort` → `choose.ErrAbort` and drop the go-fzf import.
- **Dependency**: `go-fzf` removed from `go.mod`.
- **Runtime**: adds a hard dependency on the `fzf` binary being on `PATH` (already required by the user's workflow); the popup needs tmux 3.3+ (falls back to full-screen outside tmux).
- **UX**: ks pickers now match the user's fzf binds and open as tmux popups; `resource_live` streams live in the popup.
- **History**: an earlier streaming attempt on the go-fzf library (and full-screen fzf) corrupted rendering in the target terminal (WSL2 + Windows Terminal + tmux); the `--tmux` popup is the fix that lets streaming render cleanly.
