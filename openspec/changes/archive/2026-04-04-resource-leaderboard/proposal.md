## Why

The vote-weighted resource selection introduced in `resource-vote-sort` accumulates usage data in `.ks.resources.json`, but there is no way to inspect that data without reading the raw JSON. A `resource leaderboard` command surfaces this data as a readable ranked table.

## What Changes

- New `resource_leaderboard` Cobra command registered in `cmd/ks/`
- New `ResourceLeaderboard` business logic function in `internal/cmds/`
- Output: formatted table printed to stdout — rank, resource name, vote count — sorted by votes descending then alphabetically, only resources with at least 1 vote shown (configurable via flag)

## Capabilities

### New Capabilities
- `resource-leaderboard`: CLI command that reads `.ks.resources.json` and prints a ranked, formatted table of resource selection counts

### Modified Capabilities
<!-- none -->

## Impact

- `cmd/ks/resource_leaderboard.go`: new Cobra command file
- `internal/cmds/resource_leaderboard_cmd.go`: business logic
- `cmd/ks/main.go`: register the new command in the `commands` slice
- No new external dependencies — uses `text/tabwriter` from stdlib for table formatting
