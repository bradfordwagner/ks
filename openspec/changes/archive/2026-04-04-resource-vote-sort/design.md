## Context

`.ks.resources.json` currently stores two top-level fields: `cache` (tmux-pane → last-used resource) and `names` (flat `[]string` of all known resource types). `resource load` completely overwrites the file by fetching from the Kubernetes API, which would erase any per-resource metadata added in the future. `resource` and `resource all` present the list alphabetically with no usage weighting.

The goal is to add a lightweight voting mechanism so the most-used resources float to the top of the fzf picker without any explicit user configuration.

## Goals / Non-Goals

**Goals:**
- Versioned JSON schema with automatic v1→v2 migration on load
- `names` entry per resource carries a `votes` integer
- Vote is incremented each time the user selects a resource via `resource` or `resource all`
- fzf list is presented sorted: descending votes, then ascending alphabetical
- `resource load` merges server-discovered names into the existing list, preserving votes

**Non-Goals:**
- Vote decay or time-weighted scoring
- Per-kubeconfig or per-namespace vote isolation
- Resetting votes via a CLI flag (can be done by editing the JSON)
- Migrating v2 back to v1

## Decisions

### D1: Version field as integer constant

Store `version` as an `int` on the `Resources` struct, JSON key `"version"`. Absent/zero → v1. This is simpler than a string semver and sufficient for the small number of schema generations anticipated.

*Alternative considered*: string enum (`"v1"`, `"v2"`). Rejected — integer comparison is cleaner and Go's zero-value handles missing field automatically.

### D2: `ResourceEntry` replaces `[]string`

```go
type ResourceEntry struct {
    Name  string `json:"name"`
    Votes int    `json:"votes"`
}
```

`Resources.Names` changes type from `[]string` to `[]ResourceEntry`. The JSON key stays `"names"` to minimise diff noise; the type change handles migration naturally (v1 string array → v2 object array).

*Alternative considered*: parallel map `map[string]int` for votes. Rejected — slice preserves insertion/sort order and serialises more readably.

### D3: Migration inside `LoadResources`

After unmarshalling, if `version < 2`, call `migrate()` which converts any legacy string-keyed data. Since v1 stored `names` as `[]string` (which will fail to unmarshal into `[]ResourceEntry`), a two-pass decode is used: first unmarshal into a shim struct with `NamesRaw json.RawMessage`, then attempt to decode as `[]ResourceEntry`; if that fails, decode as `[]string` and synthesize entries with `votes: 0`.

*Alternative considered*: separate migration binary. Rejected — transparent in-process migration is the simplest UX.

### D4: `resource load` merges, not replaces

Load the existing `.ks.resources.json` first (if present). Build a map of `name → votes` from the existing entries. For each name returned by the Kubernetes API, upsert into the map. Write back as v2. This preserves votes accumulated before a `resource load`.

*Alternative considered*: always reset votes on load. Rejected — defeats the purpose of vote accumulation.

### D5: Vote increment happens before k9s launch

`Resource()` increments the vote and writes the file (async goroutine, same pattern as `Upsert`) before calling `k9s.Run`. Order: select → upsert pane cache → increment vote → write → launch k9s.

## Risks / Trade-offs

- **Concurrent writes** → Mitigation: writes are already fire-and-forget goroutines; race is benign (last write wins) and rare in practice.
- **Large vote counts** → Mitigation: `int` (64-bit on modern platforms) is sufficient for any realistic usage lifetime.
- **v1 JSON with non-standard `names` content** → Mitigation: two-pass decode; if second pass also fails, treat as empty list and log a warning.

## Migration Plan

1. User upgrades `ks` binary.
2. Next invocation of any command that calls `LoadResources` triggers in-place v1→v2 migration.
3. File is written back as v2 on the next write (vote increment or explicit `resource load`).
4. No manual steps required; no rollback needed (old binary will error on the new schema's `names` field type, but the user can delete the file to reset).
