## Context

`internal/choose.One`/`Multi` wrap `github.com/koki-develop/go-fzf@v0.15.0` and are called from ~9 sites across `cmd/ks/main.go` and `internal/cmds/*.go`. The module's `option.go` was read directly (not from docs) to confirm its `Option` surface: `WithLimit`, `WithNoLimit`, `WithPrompt`, `WithCursor`, `WithSelectedPrefix`, `WithUnselectedPrefix`, `WithStyles`, `WithKeyMap`, `WithInputPlaceholder`, `WithCountViewEnabled`, `WithCountView`, `WithHotReload`, `WithCaseSensitive`, `WithInputPosition`. There is no header/title widget — `WithPrompt` (default `"> "`) is the only option that renders persistent text next to the input line (`model.go`: `input.Prompt = opt.prompt`). `styles.go` was also read directly and confirmed `WithStyles(fzf.WithStylePrompt(fzf.Style{...}))` applies a `lipgloss.Style` to the whole prompt string (`input.PromptStyle = opt.styles.option.prompt`), which is how the prompt's color is set.

## Goals / Non-Goals

**Goals:**
- Show the active kubeconfig's current-context name as a prefix on every fzf prompt in the app, with graceful fallback when it can't be determined.
- Color that prefix Catppuccin Mocha blue (`#89b4fa`) so it stands out from the rest of the picker UI.

**Non-Goals:**
- No real "header line" above the list — the go-fzf version in use doesn't support one; the prompt-prefix approach is the ceiling of what's achievable without vendoring/forking the library or swapping libraries.
- No per-substring styling (e.g. brackets in one color, context name in another) — go-fzf's `WithStylePrompt` applies one `lipgloss.Style` to the entire prompt string, so `"[<context>] > "` is colored as a single unit.
- No support for other Catppuccin flavors (Latte/Frappé/Macchiato) or user-configurable color — Mocha blue is hardcoded as the flagship/default flavor; revisit if the user wants a different flavor or a `--color`/env override later.
- No caching/memoization of `CurrentContext` — it's a cheap local YAML parse called once per picker invocation.
- Kubeconfig-listing filtering (junk files in `internal/list.Kubeconfigs`) is out of scope for this change — it was implemented and then explicitly walked back.

## Decisions

- **`fzf.WithPrompt` for the context indicator**: confirmed via source read as the only available persistent-text hook; formatted as `"[<context>] > "` to keep the familiar `"> "` cursor cue while prefixing context.
- **`fzf.WithStyles(fzf.WithStylePrompt(...))` for color, hardcoded Catppuccin Mocha blue `#89b4fa`**: the only styling hook that targets the prompt specifically (as opposed to matches, cursor, etc.); applied as a small package-level `promptStyle()` helper in `internal/choose/choose.go` and added to both `One` and `Multi`'s option list alongside `WithPrompt`.
- **`choose.One`/`choose.Multi` take a required `header string` parameter** (not optional/variadic): every call site has a meaningful kubeconfig in scope (`a.Kubeconfig`), so making it required keeps the indicator consistent everywhere rather than silently omitted at a forgotten call site.
- **`kube.CurrentContext` falls back to `filepath.Base(kubeconfig)`** rather than erroring, matching the existing tolerance pattern in `kube.SetNamespace`/`kube_new_ns_cmd.ResolveIfSymlink` — a missing/malformed kubeconfig shouldn't block showing a picker.
- **Root menu loads `standardArgs` via `flag_helper.Load` before its own picker**: currently only subcommand `RunE` functions call this; the root command's own `RunE` (in `cmd/ks/main.go`) never populates `standardArgs`, so it must load it once to get `a.Kubeconfig` for the header. Subcommands reload it again themselves afterward (existing behavior, harmless).
- **`resource_cmd.go`'s `resolveResourceType` gains a `header string` parameter**, used only on its fzf-fallback branch (env-var override and pane-cache hits skip fzf entirely and ignore it).

## Risks / Trade-offs

- [`fzf.WithPrompt` prefix is less visually distinct than a real header would be] → Accepted given library constraints; still a strict improvement over no indicator at all.
- [Changing `choose.One`/`choose.Multi` signatures is a breaking internal API change] → Contained entirely within this repo; every call site is being updated in the same change, and `go build ./...` will catch any missed site.

## Open Questions

None outstanding — go-fzf's option surface was directly verified against the vendored source rather than assumed.
