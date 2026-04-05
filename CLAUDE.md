# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

**Install/Build:**
```bash
go install ./cmd/ks
```

**Test:**
```bash
go test ./...
go test ./internal/resources/...  # single package
```

**Lint/Vet:**
```bash
go vet ./...
```

**Cross-platform release build:**
```bash
goreleaser build --snapshot --clean
```

## Architecture

`ks` is a Kubernetes CLI helper that wraps kubeconfig management, tmux pane orchestration, and k9s (terminal UI) launching into interactive workflows using fzf for selection.

### Entry Point Flow

`cmd/ks/main.go` registers ~12 Cobra commands. When invoked without arguments, the root command presents an fzf menu of commands to run (a recursive self-selection loop). Each `cmd/ks/*.go` file wires flags and delegates to `internal/cmds/`.

### Package Responsibilities

| Package | Role |
|---------|------|
| `internal/args` | Shared flag struct: `Dir`, `Kubeconfig`, `Timeout` |
| `internal/choose` | fzf wrapper — `One()` (single pick) and `Multi()` (multi-select) |
| `internal/cmds` | Business logic for each command |
| `internal/kube` | k8s clientset init, kubeconfig loading, namespace setting |
| `internal/k9s` | Launches k9s TUI (headless-capable) |
| `internal/tmux` | Splits panes, creates windows, loads KUBECONFIG into tmux buffer |
| `internal/resources` | JSON cache (`.ks.resources.json`) keyed by tmux session+pane |
| `internal/link` | Symlink management for `~/.kube/config` |
| `internal/list` | Lists kubeconfig files from a directory |

### Resource Caching

`internal/resources/types.go` maintains a per-tmux-pane JSON cache of the last-used Kubernetes resource type. The cache key is derived from the `$TMUX` environment variable (session + pane ID). Commands like `resource` read this cache to default to the previously selected resource, then upsert on selection.

### Tmux Integration Pattern

Commands that spawn tmux panes/windows set `KUBECONFIG=<path>` in the new pane's environment. `tmux_multi` uses multi-select fzf to pick multiple kubeconfigs and opens one pane per selection in tiled layout.

### Key Dependencies

- `github.com/spf13/cobra` — CLI framework
- `github.com/koki-develop/go-fzf` — interactive fuzzy selection
- `k8s.io/client-go` — Kubernetes API (namespace listing, kubeconfig mutation)
- `github.com/bradfordwagner/go-util` — shared utilities

### CI

GitHub Actions workflows in `.github/workflows/` delegate to external taskfiles (`github.com/bradfordwagner/taskfiles`). Branch pushes trigger a snapshot build; git tags trigger a full GoReleaser release with checksums.
