# ks

- [y] ks pipe - KUBECONFIG=$(ks pipe)
- [y] ks tmux_multi
- [y] ks tmux_window
  - [y] remove kcc hack
- [y] ks link
- [y] resources load
- existing commands
  - [y] resource
  - [x] ns_resource
  - [y] select_ns
  - [y] new_ctx_ns_cp (kube_new_ns)
  - [y] select_primary_kube_ctx
  - [y] select_kube_ctx_cp (kube_cp)
  - [y] multi_kube_ctx
  - [y] kube_new_window
  - [y] allns_resource (all_resource)
- [y] jump off to all cmds
  - [y] allow completion ordering for default


## tmux-resurrect Integration

`ks save` and `ks restore` integrate with [tmux-resurrect](https://github.com/tmux-plugins/tmux-resurrect) to restore `KUBECONFIG` and the active k9s resource view in each pane after a session resurrect.

### How it works

- **`ks save`** (pre-save hook): scans all tmux panes, extracts `KUBECONFIG` from each pane's process environment via `/proc`, cross-references the resource cache, and writes a positional sidecar `~/.ks/.ks.resurrect.json` keyed by `session/window/pane` index (stable across resurrect)
- **`ks restore`** (post-restore hook): reads the sidecar, matches panes by positional index, and sends `KS_RESOURCE=<resource> KUBECONFIG=<path> ks resource` (or `ks resource_all`) inline into each pane — no fzf, no manual re-selection

### Setup

Add to `~/.tmux.conf`:

```tmux
set -g @resurrect-hook-pre-save 'ks save'
set -g @resurrect-hook-post-restore-all 'ks restore'
```

### Notes

- `ks save` requires Linux (`/proc` filesystem); it is a no-op on macOS
- Panes with neither a `KUBECONFIG` nor a cached resource are silently skipped
- Panes with only a resource (no `KUBECONFIG`) have the resource and verb restored without setting `KUBECONFIG`
- Panes with `KUBECONFIG` but no cached resource have only `KUBECONFIG` restored; resource selection falls back to fzf on next `ks resource` invocation
- Verb (`resource` vs `resource_all`) is detected from the running `k9s` or `ks` process cmdline at save time
- `ks clear_pane` and `ks clear_cache` work normally after restore — `KS_RESOURCE` is never exported into the shell environment

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `KS_DIR` | `~/.kube` | Directory scanned for kubeconfig files |
| `KS_DATA_DIR` | `~/.ks` | Directory for ks data files (resource cache, resurrect sidecar) |
| `KUBECONFIG` | `~/.kube/config` | Active kubeconfig path (standard k8s var) |
| `KS_RESOURCE` | — | Single-invocation resource override for `ks resource` / `ks resource_all`; bypasses fzf and cache lookup without being exported into the shell |

## Local
```bash
go install ./cmd/ks
```
