# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

**Install/Build:**
```bash
task        # go install ./cmd/ks
go install ./cmd/ks  # equivalent, without task
```

**Test:**
```bash
go test ./...
go test ./internal/resources/...          # single package
go test ./internal/resources/... -run TestResources/SortedNames  # single spec
```

Tests use [Ginkgo v2](https://onsi.github.io/ginkgo/) + Gomega (BDD-style). Each package with tests has a `suite_test.go` that bootstraps `RunSpecs`. Use `-run <suite>/<describe>/<it-text>` to target a single spec; Ginkgo also supports `FIt`/`FDescribe` focus markers.

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

`cmd/ks/main.go` registers all Cobra commands. Without arguments, the root command presents an fzf menu and loops: pick a command, run it, repeat (exit via Esc/Ctrl-C). Each `cmd/ks/*.go` file wires flags and delegates to `internal/cmds/`.

All commands share the same `standardArgs` (type `args.Standard`) with four flags:

| Flag | Short | Env | Default |
|------|-------|-----|---------|
| `--dir` | `-d` | `KS_DIR` | `~/.kube` |
| `--data-dir` | — | `KS_DATA_DIR` | `~/.ks` |
| `--kubeconfig` | `-k` | `KUBECONFIG` | `~/.kube/config` |
| `--timeout` | `-t` | — | `10s` |

`--dir` is the kubeconfig discovery directory. `--data-dir` is where ks writes its own data files (resource cache, resurrect sidecar).

Flags are populated from env vars via `flag_helper.Load` at command runtime (not at init time).

### Command Reference

| Command | What it does |
|---------|-------------|
| `resource` | fzf-pick a k8s resource type → launch k9s at that view; skips fzf if pane already has a cached resource |
| `resource_all` | same as `resource` but passes `-A` (all namespaces) to k9s |
| `resource_load` | calls the cluster's discovery API, writes resource names to `~/.ks/.ks.resources.json`; **must run before `resource`** |
| `resource_leaderboard` | tabular view of resource usage counts; `--all` includes zero-vote entries |
| `set_ns` | fzf-pick a namespace → mutates the kubeconfig's current-context namespace |
| `tmux_multi` | multi-select kubeconfigs → split one tmux pane per selection, each with `KUBECONFIG` set |
| `tmux_window` | fzf-pick one kubeconfig → open a new tmux window with `KUBECONFIG` set |
| `kube_cp` | fzf-pick a kubeconfig → load its path into the tmux clipboard buffer |
| `kube_new_ns` | fzf-pick a kubeconfig → create a new namespace in that cluster |
| `link` | fzf-pick a kubeconfig → symlink it to `~/.kube/config` |
| `pipe` | fzf-pick a kubeconfig → print its full path to stdout (no newline); designed for shell pipes |
| `clear_cache` | wipe the entire tmux→pane→resource cache in `~/.ks/.ks.resources.json` |
| `clear_pane` | remove only the current tmux pane's cached resource entry |
| `save` | tmux-resurrect pre-save hook: snapshot KUBECONFIG + resource per pane to `~/.ks/.ks.resurrect.json` |
| `restore` | tmux-resurrect post-restore hook: re-apply KUBECONFIG and re-launch `ks resource`/`ks resource_all` per pane |

### Package Responsibilities

| Package | Role |
|---------|------|
| `internal/args` | Shared flag struct: `Dir`, `DataDir`, `Kubeconfig`, `Timeout` |
| `internal/choose` | fzf wrapper — `One()` (single pick) and `Multi()` (multi-select) |
| `internal/cmds` | Business logic for each command |
| `internal/kube` | k8s clientset init, kubeconfig loading, namespace setting |
| `internal/k9s` | Launches k9s TUI (headless-capable) |
| `internal/tmux` | Splits panes, creates windows, loads KUBECONFIG into tmux buffer |
| `internal/resources` | JSON cache (`~/.ks/.ks.resources.json`) keyed by tmux session+pane |
| `internal/resurrect` | Sidecar types, `/proc` env reading, tmux pane listing for save/restore |
| `internal/link` | Symlink management for `~/.kube/config` |
| `internal/list` | Lists kubeconfig files from a directory |

### Resource Cache (`internal/resources`)

`~/.ks/.ks.resources.json` (v2 schema) stores two things:
- `names`: `[]ResourceEntry{Name, Votes}` — every known k8s resource type with selection counts
- `cache`: nested map `$TMUX → $TMUX_PANE → resource-name` — the last resource selected per pane

On `resource` invocation: if the current pane already has a cached resource it skips fzf, increments that resource's vote, and writes the file asynchronously (`go r.Write(...)`) so k9s starts without waiting.

`resource_load` preserves existing votes when refreshing — it merges discovery results with the existing file. The file migrates automatically from v1 (flat `[]string` names) to v2 (`[]ResourceEntry`) on first load.

### Tmux Integration Pattern

Commands that spawn panes/windows set `KUBECONFIG=<path>` in the new pane's environment. `tmux_multi` uses multi-select fzf to pick multiple kubeconfigs and opens one pane per selection in tiled layout.

### tmux-resurrect Wiring

Add to `~/.tmux.conf` to enable session restore:

```tmux
set -g @resurrect-hook-pre-save 'ks save'
set -g @resurrect-hook-post-restore-all 'ks restore'
```

`ks save` reads `/proc/<pid>/environ` (Linux only) to extract `KUBECONFIG` and cross-references the resource cache. `ks restore` sends `KS_RESOURCE=<r> KUBECONFIG=<path> ks <verb>` inline into each matched pane — `KS_RESOURCE` is not exported so `clear_pane`/`clear_cache` work normally afterward.

### CI

GitHub Actions workflows in `.github/workflows/` delegate to external taskfiles (`github.com/bradfordwagner/taskfiles`). Branch pushes trigger a snapshot build; git tags trigger a full GoReleaser release with checksums. `workflow.yaml` at the repo root is an Argo Workflows template used for alternative CI builds via `quay.io/bradfordwagner/go-builder`.
