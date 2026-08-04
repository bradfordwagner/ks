## Context

See proposal.md. Building blocks already present: `internal/tmux` splits a tiled pane with `KUBECONFIG` set and sends keys to it; `resource_cmd.go` launches k9s for one resource and honors the `KS_RESOURCE` override (used by split panes and resurrect restore); `internal/resources` stores vote-sorted names + a per-pane cache; `kube.Dynamic` gives a dynamic client. This change composes those rather than adding parallel machinery.

## Goals / Non-Goals

**Goals:**
- One selection model (`resource`) covering single and multi; single-pick unchanged.
- Reuse the proven `KS_RESOURCE=<r> ks <verb>` pane-launch mechanism for splits and for restore.
- Cheap, scalable usage detection (existence, not counts) with a cached and a live entry point.

**Non-Goals:**
- Exact object counts (existence only — chosen for cost).
- The picker engine itself (fzf binary, `FZF_DEFAULT_OPTS`, streaming primitive) — that is the `interactive-picker` capability; this change only consumes it.
- Per-resource namespace scope; an in-picker filter hotkey.

## Decisions

### D1: Fold multi-select into `resource`; remove `resource_multi`

`Resource` resolves a single resource from `KS_RESOURCE`/cache first (fast path → current pane). Otherwise it multi-selects: `launchSelected` opens split panes for `selected[1:]` (each `KS_RESOURCE=<r> ks resource[_all]`) then runs k9s for `selected[0]` in the current pane (splits must happen before the blocking `k9s.Run`).
- **Why:** single unified verb; single-pick users see no change; split children reuse the existing override path so caching/voting/restore all work unchanged.
- **Alternative:** keep a separate `resource_multi` (the prior change). Rejected — redundant once `resource` multi-selects.

### D2: Existence probe, all namespaces, bounded parallelism

`probePopulated` runs a `Limit:1` list per GVR via the dynamic client, cluster-wide (no namespace → spans all namespaces for namespaced types), capped at `probeConcurrency` and bounded by a `context.WithTimeout(--timeout)`. It returns a map keyed only by successfully-checked names (value = has ≥1), so callers distinguish known-empty from unknown. `listableGVRs` builds the name→GVR map from discovery (listable, non-subresource; first version wins), tolerating partial discovery.
- **Why:** existence is the only scalable signal across a whole catalog; all-namespaces matches "is this type in use at all".
- **Alternative:** exact counts (list everything) — rejected as too expensive.

### D3: Cached (`resource_exists`) vs live (`resource_live`)

`resource_load` persists `Populated`/`Probed` per entry (schema v3, additive). `resource_exists` reads those flags — instant, no cluster calls. `resource_live` probes fresh each run. Both funnel into `launchSelected`, so selection/splits behave identically to `resource`.
- **Why:** a fast cached path for the common case and an always-fresh path when the cache is stale; the ladder is resource (all) → resource_exists (cached populated) → resource_live (live populated).

### D4: `resource_live` streams into the tmux popup

`resource_live` dispatches probes in `prioritizeCandidates` order and feeds each populated type into `choose.MultiStream` over a channel; fzf reads its input incrementally, so types appear as the scan finds them, rendered in the tmux popup (see the interactive-picker capability).
- **History:** streaming under full-screen fzf (and earlier the go-fzf library) corrupted rendering in the target terminal (WSL2 + Windows Terminal + tmux) — overlapping/uncleared text, ghost info lines. Moving the finder into a tmux popup, where tmux owns the screen region, is what makes incremental redraws render cleanly.
- **Trade-off:** results stream in probe-arrival order (dispatch is priority-first, so it approximates the sort) rather than a strict post-hoc sort. `resource_exists` remains the instant, no-scan path.

### D5: `resource_live` non-persistence of flags

`resource_live` writes only votes + pane cache on selection (via the shared `runSingle`); it never rewrites `Populated`/`Probed` — only `resource_load` does.
- **Why:** keeps the flag a deliberate load-time snapshot rather than a side effect of watching.

## Risks / Trade-offs

- **Probe latency / API load on large catalogs.** → Bounded concurrency + `--timeout`; failures non-fatal. `resource_live` streams so results appear as found; `resource_exists` avoids the live cost entirely.
- **Streamed list is in probe-arrival order, not a strict sort.** → Accepted (see D4); dispatch is priority-first so it approximates the sort.
- **Streaming rendering depends on the tmux popup.** → Full-screen fzf corrupts on incremental redraws in the target terminal; the popup (interactive-picker capability) is required for clean streaming. Outside tmux there is no popup, so streaming there is best-effort.
- **Stale `resource_exists` results.** → By design (cached snapshot); `resource_live` is the fresh alternative.
- **Probe needs the dynamic client / RBAC to list.** → Forbidden types are simply left unprobed, never hidden as "known empty".
