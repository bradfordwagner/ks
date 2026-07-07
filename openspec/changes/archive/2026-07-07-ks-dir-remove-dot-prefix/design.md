## Context

`ks` stores two runtime data files under `~/.ks` (the `DataDir`):
- `.ks.resources.json` — k8s resource cache and vote leaderboard
- `.ks.resurrect.json` — tmux-resurrect session snapshot

Both filenames begin with `.`, making them hidden from standard `ls` output and file browsers. The `~/.ks` directory is already dedicated to `ks` data, so there is no need to hide individual files within it.

The filenames are defined as package-level constants:
- `internal/resources/types.go` → `CacheFile = ".ks.resources.json"`
- `internal/resurrect/types.go` → `ResurrectFile = ".ks.resurrect.json"`

All callers pass `DataDir` to `Load`/`Write` functions that join with these constants; no caller hardcodes the filename directly.

## Goals / Non-Goals

**Goals:**
- Remove the dot prefix from both data file constants so files are visible by default
- Update all inline documentation referencing the old names

**Non-Goals:**
- Auto-migrating existing files on disk (out of scope; users re-run `ks resource_load` and `ks save`)
- Changing the `DataDir` location (`~/.ks`) itself
- Changing the JSON schema or data format

## Decisions

**Change constants only, not callers**: All file I/O goes through `resources.LoadResources(DataDir)` / `r.Write(DataDir)` and `resurrect.Load(DataDir)` / `state.Write(DataDir)`. Updating the two constants is sufficient — no caller changes needed.

**No migration code**: The renamed files are regenerated cheaply (`ks resource_load` rebuilds the resource cache; `ks save` rebuilds the resurrect snapshot). Adding migration logic would complicate the binary for a one-time user action.

**Simpler names**: Drop the redundant `ks.` infix since the files already live under `~/.ks`:
- `.ks.resources.json` → `resources.json`
- `.ks.resurrect.json` → `resurrect.json`

## Risks / Trade-offs

- **Existing files become stale** → Users with existing `~/.ks/.ks.*.json` files will lose their vote history and resurrect snapshot after upgrading. Mitigation: document in CLAUDE.md and release notes; vote counts regenerate organically with use.
- **No rollback path** → If a user downgrades, the old binary will not find the new filenames. Mitigation: low risk given single-user personal tooling.
