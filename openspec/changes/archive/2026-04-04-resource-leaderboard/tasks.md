## 1. Business Logic

- [x] 1.1 Create `internal/cmds/resource_leaderboard_cmd.go` with `ResourceLeaderboard(a args.Standard, all bool)` that loads `.ks.resources.json` and handles missing-file case with a friendly message
- [x] 1.2 Filter entries to votes ≥ 1 unless `all` is true; if nothing to show, print "no usage data" message and return
- [x] 1.3 Sort entries by votes descending then name ascending (reuse `SortedNames` order logic directly on `[]ResourceEntry`)
- [x] 1.4 Render ranked table to stdout using `text/tabwriter` with columns `#`, `RESOURCE`, `VOTES`

## 2. Cobra Command & Registration

- [x] 2.1 Create `cmd/ks/resource_leaderboard.go` with Cobra command `resource_leaderboard`, standard flags, and `--all` bool flag
- [x] 2.2 Register `resourceLeaderboardCmd` in the `commands` slice in `cmd/ks/main.go`

## 3. Tests

- [x] 3.1 Unit test missing file: `ResourceLeaderboard` prints message and returns nil when no cache file exists
- [x] 3.2 Unit test default filter: zero-vote entries are excluded; only voted resources appear in output
- [x] 3.3 Unit test `--all` flag: zero-vote entries appear when `all=true`
- [x] 3.4 Unit test sort order: votes desc then alpha asc across multiple entries
- [x] 3.5 Unit test table format: header row present, each data row contains rank, name, and vote count
