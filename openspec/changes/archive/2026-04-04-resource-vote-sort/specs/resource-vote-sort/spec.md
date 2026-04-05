## ADDED Requirements

### Requirement: Schema versioning
The `Resources` struct SHALL include an integer `version` field persisted as `"version"` in JSON. A missing or zero value SHALL be treated as version 1 (v1). The current supported version is 2 (v2).

#### Scenario: New file written as v2
- **WHEN** `LoadResources` is called and no file exists, then any write occurs
- **THEN** the resulting JSON SHALL contain `"version": 2`

#### Scenario: Existing v1 file is loaded
- **WHEN** `LoadResources` reads a file with no `"version"` field
- **THEN** the returned `Resources` struct SHALL have `Version == 2` after migration

### Requirement: ResourceEntry type
The `names` list SHALL be stored as an array of objects, each containing a `name` string and a `votes` integer. The JSON key SHALL remain `"names"`.

#### Scenario: Serialised format
- **WHEN** a `Resources` value is marshalled to JSON
- **THEN** `"names"` SHALL be an array of `{"name":"<resource>","votes":<int>}` objects

### Requirement: V1 to V2 migration
`LoadResources` SHALL automatically migrate v1 files (where `names` is a `[]string`) to v2 format, initialising each entry's `votes` to 0.

#### Scenario: V1 names preserved after migration
- **WHEN** a v1 file with `"names":["pods","services"]` is loaded
- **THEN** the v2 result SHALL contain entries for `pods` and `services` each with `votes: 0`

#### Scenario: Migration is idempotent
- **WHEN** a v2 file is loaded and immediately written back
- **THEN** the resulting file SHALL be identical in schema (no double-migration)

### Requirement: Vote increment on resource selection
`Resource()` and `Resource(all=true)` SHALL increment the `votes` counter for the selected resource entry after the user makes a selection, then persist the updated file asynchronously.

#### Scenario: Vote count increases on selection
- **WHEN** a user selects "pods" via the `resource` command
- **THEN** the `votes` field for `pods` in `.ks.resources.json` SHALL be incremented by 1

#### Scenario: Vote increment on resource all
- **WHEN** a user selects "deployments" via the `resource all` command
- **THEN** the `votes` field for `deployments` SHALL be incremented by 1

### Requirement: Vote-weighted sort for fzf display
The resource list presented to fzf SHALL be sorted by `votes` descending, then by `name` ascending for entries with equal votes.

#### Scenario: Higher-voted resource appears first
- **WHEN** "pods" has 5 votes and "services" has 2 votes
- **THEN** fzf SHALL display "pods" before "services"

#### Scenario: Alphabetical tiebreak within same vote count
- **WHEN** "configmaps" and "pods" both have 3 votes
- **THEN** fzf SHALL display "configmaps" before "pods"

#### Scenario: Zero-vote resources sorted alphabetically
- **WHEN** all resources have 0 votes
- **THEN** fzf SHALL display resources in alphabetical order

### Requirement: Resource load preserves votes
`ResourceLoad` SHALL merge Kubernetes API-discovered resource names into the existing `.ks.resources.json`, preserving vote counts for resources already present in the file.

#### Scenario: Existing votes retained after load
- **WHEN** "pods" has 5 votes in `.ks.resources.json` and `resource load` is run
- **THEN** "pods" SHALL still have 5 votes in the updated file

#### Scenario: New resources added with zero votes
- **WHEN** `resource load` discovers a resource name not present in the current file
- **THEN** that resource SHALL be added with `votes: 0`

#### Scenario: Removed resources are pruned
- **WHEN** a resource name in `.ks.resources.json` is no longer returned by the Kubernetes API
- **THEN** that entry SHALL be removed from the file after `resource load`
