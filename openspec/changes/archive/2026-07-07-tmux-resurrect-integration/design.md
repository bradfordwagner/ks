## Context

`ks` stores per-pane k8s state in `~/.kube/.ks.resources.json` keyed by `$TMUX` + `$TMUX_PANE`. Both are volatile: tmux-resurrect assigns new IDs on every restore, so all cache entries become stale and `KUBECONFIG` (set via `-e` at pane creation) is lost.

tmux-resurrect reconstructs sessions positionally — `session_name + window_index + pane_index` are stable across save/restore. This gives us a stable key. Linux's `/proc/<pid>/environ` lets us read a process's inherited environment at any time, solving the KUBECONFIG retrieval problem without modifying pane creation code.

## Goals / Non-Goals

**Goals:**
- Restore `KUBECONFIG` and re-launch `ks resource` / `ks resource_all` in each ks-managed pane after resurrect, without fzf interaction
- Zero changes to existing pane creation commands (`tmux_window`, `tmux_multi`)
- `clear_pane` / `clear_cache` semantics unchanged after restore

**Non-Goals:**
- macOS support for `ks save` (no `/proc` filesystem)
- Automatic tmux.conf hook setup (user wires `@resurrect-hook-pre-save` / `@resurrect-hook-post-restore-all` manually)
- Restoring panes not managed by `ks` (no `KUBECONFIG` → silently skipped)
- Restoring arbitrary shell state beyond KUBECONFIG + resource

## Decisions

### 1. `/proc` for env extraction vs. proactive sidecar at creation time

**Decision:** `/proc`

Proactive tracking would require modifying `Split` and `NewWindow` to write positional data at creation time — touching more code and adding I/O to the hot path. `/proc/<pid>/environ` reads the environment of any running process retroactively, requiring no changes to existing commands. The only constraint is Linux.

### 2. Inline `KS_RESOURCE=<r> ks <verb>` vs. external cache pre-population

**Decision:** Inline env var (`KS_RESOURCE`)

External pre-population requires the restore command to know the new `$TMUX_PANE` ID from outside the pane, then write a cache entry, then send `ks resource`. The inline approach is cleaner: restore sends `KS_RESOURCE=pods ks resource` into the pane; `Resource()` reads the env var, calls `Upsert` + `Write`, and from that point the pane behaves normally. After k9s exits, `KS_RESOURCE` is gone (never exported), so `clear_pane` / `clear_cache` work without interference.

### 3. Sidecar file vs. augmenting `.ks.resources.json`

**Decision:** Separate sidecar `~/.kube/.ks.resurrect.json`

`.ks.resources.json` is keyed by volatile `$TMUX`/`$TMUX_PANE` IDs; the sidecar is keyed by stable positional indices. Mixing both schemes in one file conflates two different key spaces. A separate file has a clear schema version and lifecycle independent of the resource cache.

### 4. k9s PID discovery for verb detection

**Decision:** Walk `/proc` for descendants of `pane_pid` with name `k9s`

`pane_pid` from tmux is the shell PID. When `ks resource` runs, the process tree is `shell → ks → k9s`. `ks` is gone once k9s exits, but while k9s is running, we can find it by scanning `/proc/*/status` for entries whose `PPid` chain leads to `pane_pid`. Reading `/proc/<k9s_pid>/cmdline` then reveals whether `-A` is present.

### 5. Package structure

**Decision:** New `internal/resurrect/` package

Keeps proc-reading, sidecar types, and tmux list-panes logic separate from `internal/cmds`. `cmds` delegates to `resurrect` the same way it delegates to `tmux`, `resources`, etc.

## Risks / Trade-offs

- **`go r.Write()` race at save time** → The async write in `Resource()` completes in milliseconds; resurrect-save is always user-triggered, never sub-second after restore. Not a practical risk.
- **`/proc` PID reuse** → If a pane's shell PID was reused by an unrelated process, we'd read the wrong env. Mitigated by cross-checking `pane_current_command` against the process name in `/proc/<pid>/comm`.
- **k9s not a direct child** → Process tree depth (shell → ks → k9s) requires descent. If the shell execs ks directly (no fork), k9s would be a direct child. Implementation must handle both depths.
- **Linux-only save** → `ks save` is a no-op on non-Linux. Restore still works anywhere (it only calls tmux and reads the sidecar). This is acceptable given WSL2 is the primary target.
- **Session/window/pane index collisions** → If the restored session has more panes than the saved session (e.g., user added panes before save), positional matching may mis-assign KUBECONFIG. Low probability; no mitigation planned.
