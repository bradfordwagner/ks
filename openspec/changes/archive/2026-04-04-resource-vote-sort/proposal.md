## Why

Resource selection in `ks` currently displays all Kubernetes resource types alphabetically, forcing users to scroll or type to find frequently-used resources. Adding a vote-based sort makes the most commonly accessed resources rise to the top automatically over time.

## What Changes

- Add a `version` field to `.ks.resources.json`; existing files without the field are treated as version 1 (v1)
- Introduce version 2 (v2) schema: the flat `names []string` is replaced by `names []ResourceEntry` where each entry carries a `name` and `votes` counter
- Automatic migration: when `ks` loads a v1 file it migrates it to v2 in-place (votes initialized to 0)
- `resource load` now merges newly discovered server resource names into the existing v2 list, preserving vote counts for resources that already exist and adding new entries with 0 votes
- `resource` and `resource all` commands increment the vote count for whichever resource type the user selects before launching k9s
- fzf selection list is sorted: highest votes first, then alphabetically within the same vote count

## Capabilities

### New Capabilities
- `resource-vote-sort`: Vote-weighted sorting of Kubernetes resource types; vote incremented on selection; persisted in `.ks.resources.json` v2 schema with auto-migration from v1

### Modified Capabilities
<!-- no existing spec files to delta -->

## Impact

- `internal/resources/types.go`: new `ResourceEntry` type, `version` field on `Resources`, migration logic, vote-increment helper, sorted-names accessor
- `internal/cmds/resource_load_cmd.go`: load existing file before overwriting so votes are preserved during merge
- `internal/cmds/resource_cmd.go`: call vote-increment + write after user selects a resource
- `.ks.resources.json` on disk: schema change (non-breaking via migration); old files are upgraded automatically on first read
