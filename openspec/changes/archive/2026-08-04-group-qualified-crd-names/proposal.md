## Why

`resource_load` discards the API group and stores only the short plural name (`randomsecrets`), so the fzf picker has no group text to match on — you cannot type `redhatcop` / `redcop` to filter down to the vault-config-operator CRDs, and same-named resources in different groups collapse into one entry. The fzf wrapper matches on the string it displays (no hidden search key), so the group must live in the name itself.

## What Changes

- `resource_load` names resources from **custom / non-standard API groups** as `<plural>.<group>` (e.g. `randomsecrets.redhatcop.redhat.io`). Core and standard Kubernetes groups (`apps`, `batch`, `autoscaling`, `extensions`, `policy`, and any `*.k8s.io` group) keep their short plural name.
- The qualified name flows unchanged through the picker, the per-pane resource cache, vote counting, and the `k9s -c` argument — one identifier everywhere.
- Existing votes are preserved when a resource's stored name shifts from its bare form to its qualified form.

## Capabilities

### New Capabilities
- `resource-discovery-naming`: how discovered resource types are named for the picker — CRDs carry their API group so they are searchable by group and unambiguous across groups.

### Modified Capabilities
<!-- None. No existing spec covers resource_load's naming behavior. -->

## Impact

- **Code**: `internal/cmds/resource_load_cmd.go` (group-qualification logic + pure helpers), new Ginkgo test.
- **Data**: on the next `resource_load`, CRD entries in `resources.json` change to the qualified form; votes are migrated so the leaderboard/sort do not reset.
- **k9s**: `k9s -c <plural>.<group>` (GVR form) must resolve — it does for group-qualified resources; to be confirmed against a live cluster.
- **Docs**: `CLAUDE.md` notes that CRD names are group-qualified and filterable by group.
- No breaking API/flag changes; no new dependencies.
