# resource-usage-filtering Specification

## Purpose
Defines how ks discovers whether resource types hold objects and lets the user narrow the picker to populated types, either from a cached load-time probe or a fresh live scan.
## Requirements
### Requirement: resource_load probes usage

`resource_load` SHALL check each listable resource type for at least one object cluster-wide (across all namespaces), in parallel and bounded by the `--timeout`, and record per type whether it was successfully checked (`probed`) and whether it had ≥1 object (`populated`). Probe failures (RBAC, timeout) SHALL be non-fatal and SHALL leave the type unprobed rather than failing the load.

#### Scenario: A populated type is recorded
- **WHEN** `resource_load` runs and `randomsecrets.redhatcop.redhat.io` has at least one object
- **THEN** its entry is recorded with `probed` true and `populated` true

#### Scenario: An empty type is recorded as probed-but-not-populated
- **WHEN** `resource_load` runs and a listable type has no objects
- **THEN** its entry is recorded with `probed` true and `populated` false

#### Scenario: A probe failure is non-fatal
- **WHEN** checking a type errors (e.g. forbidden by RBAC)
- **THEN** `resource_load` still completes and that type is left unprobed

### Requirement: resource_exists filters to cached populated types

`resource_exists` SHALL show, without any live cluster calls, only the resource types recorded as `probed` and `populated` at the last `resource_load`, and SHALL launch the selection using the same behavior as `resource` (single pick in the current pane, extra picks in tiled split panes).

#### Scenario: Only cached-populated types are offered
- **WHEN** the user runs `resource_exists`
- **THEN** the picker lists only types whose last-recorded state is `probed` and `populated`, and makes no live cluster calls

#### Scenario: Guidance when nothing is recorded
- **WHEN** `resource_exists` runs and no type is recorded as populated
- **THEN** the user is told to run `resource_load` and no picker is shown

### Requirement: resource_live filters to currently-populated types

`resource_live` SHALL scan the cluster live (parallel existence checks across all namespaces, bounded by the `--timeout`), offer only types that currently have ≥1 object, and launch the selection with the same behavior as `resource`.

#### Scenario: Only currently-populated types are offered
- **WHEN** the user runs `resource_live`
- **THEN** the picker offers only types that the live scan found to have at least one object

### Requirement: resource_live streams results and probes used types first

The `resource_live` picker SHALL appear immediately and populate live as the scan finds objects, rather than waiting for the whole scan (the picker renders in a tmux popup so streaming updates render cleanly — see the interactive-picker capability). Probes SHALL be dispatched so that types which were populated at the last `resource_load` AND have selection votes are checked first (votes descending, then name), so the most relevant types tend to surface first.

#### Scenario: Results appear as they are found
- **WHEN** the user runs `resource_live` against a cluster with many types
- **THEN** the picker opens right away and populated types appear in it as the scan finds them

#### Scenario: Voted, previously-populated types are checked first
- **WHEN** the scan runs
- **THEN** types that were populated at the last `resource_load` and have votes are probed before the rest, so they tend to appear first

### Requirement: resource_live does not rewrite usage flags

On selection, `resource_live` SHALL update `resources.json` the same way `resource` does (pane cache and vote count for the chosen type) but SHALL NOT overwrite the stored `populated`/`probed` flags — only `resource_load` refreshes those.

#### Scenario: Live scan results are not persisted as flags
- **WHEN** `resource_live` selects a type and k9s exits
- **THEN** the chosen type's vote count is incremented and the pane cache is updated, while every type's `populated`/`probed` flags remain as of the last `resource_load`

### Requirement: All-namespaces via the resource_live_all verb

`resource_live_all` SHALL behave like `resource_live` — live-scanning the cluster and offering only currently-populated types, streamed into the picker as found — but SHALL watch every selected type across all namespaces, equivalent to running `resource_live -A`.

#### Scenario: resource_live_all watches all namespaces for the selection
- **WHEN** the user runs `resource_live_all` and selects a type
- **THEN** the live scan behaves as it does for `resource_live`, and the selected type opens k9s watching that type across all namespaces

#### Scenario: resource_live_all is reachable from the interactive menu
- **WHEN** the user runs bare `ks` and picks `resource_live_all` from the command menu
- **THEN** it runs the same all-namespaces live scan as `ks resource_live -A`, without requiring any flag to be passed

