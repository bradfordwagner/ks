## Context

`internal/resources/types.go` already provides `SortedNames()` (votes desc, alpha asc) and the full `[]ResourceEntry` slice. The leaderboard command simply needs to read that data and render it. No new data model changes are required.

## Goals / Non-Goals

**Goals:**
- `ks resource-leaderboard` prints a ranked table to stdout
- Columns: rank (`#`), resource name, vote count
- Sorted: votes descending, then alphabetically
- Only entries with ≥1 vote shown by default; `--all` flag shows zero-vote entries too
- Column widths auto-fit to content using `text/tabwriter`

**Non-Goals:**
- Persistent storage changes
- Interactive/fzf selection
- Color/ANSI output
- Pagination

## Decisions

### D1: `text/tabwriter` for formatting

stdlib `text/tabwriter` aligns tab-separated columns without any external dependency. Output is clean in any terminal width.

*Alternative considered*: third-party table library (e.g., `tablewriter`). Rejected — adds a dependency for trivial formatting; `tabwriter` is sufficient.

### D2: Default hides zero-vote entries

A fresh install with no usage would show hundreds of resources at 0 votes, which is noise. Default filters to `votes >= 1`. `--all` flag overrides.

*Alternative considered*: always show all. Rejected — the leaderboard is most useful as a signal of actual usage, not a full resource list (that's what `resource` command is for).

### D3: Rank column is positional (1-based)

Rank is derived from position in the sorted slice, not stored. Ties in vote count share no special treatment — rank increments normally (no "tied 3rd" display).

### D4: Business logic in `internal/cmds/resource_leaderboard_cmd.go`

Follows the existing pattern: thin Cobra command in `cmd/ks/`, all logic in `internal/cmds/`.

## Risks / Trade-offs

- **Empty file / no votes**: If `.ks.resources.json` doesn't exist or all entries are zero, default mode prints a short "no data" message rather than an empty table. Mitigates confusion on first run.
- **Very long resource names**: `tabwriter` handles this naturally by expanding the name column.
