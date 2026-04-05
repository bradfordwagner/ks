## ADDED Requirements

### Requirement: Leaderboard command exists
A `resource-leaderboard` Cobra command SHALL be registered in `ks` and appear in the root fzf command picker.

#### Scenario: Command is invokable
- **WHEN** the user runs `ks resource_leaderboard`
- **THEN** the command SHALL execute without error and print output to stdout

### Requirement: Table output format
The command SHALL print a tab-aligned table with a header row and one row per displayed resource entry. Columns SHALL be: rank (`#`), resource name (`RESOURCE`), vote count (`VOTES`).

#### Scenario: Header row present
- **WHEN** there is at least one entry to display
- **THEN** the first line of output SHALL be the header `#  RESOURCE  VOTES`

#### Scenario: Each row contains rank, name, votes
- **WHEN** the leaderboard is rendered
- **THEN** each data row SHALL contain the 1-based rank, the resource name, and the integer vote count separated by tabs

### Requirement: Sort order
Rows SHALL be sorted by votes descending, then by resource name ascending for entries with equal vote counts.

#### Scenario: Higher-voted resource ranks first
- **WHEN** "pods" has 10 votes and "services" has 3 votes
- **THEN** "pods" SHALL appear in row 1 and "services" in row 2

#### Scenario: Alphabetical tiebreak on equal votes
- **WHEN** "configmaps" and "secrets" both have 5 votes
- **THEN** "configmaps" SHALL appear before "secrets"

### Requirement: Default hides zero-vote entries
By default the command SHALL only display resources with at least 1 vote.

#### Scenario: Zero-vote entries hidden by default
- **WHEN** "nodes" has 0 votes and no `--all` flag is passed
- **THEN** "nodes" SHALL NOT appear in the output

#### Scenario: No data message when all votes are zero
- **WHEN** no resources have any votes and `--all` is not passed
- **THEN** the command SHALL print a message indicating no usage data is available and exit cleanly

### Requirement: --all flag shows zero-vote entries
When `--all` is passed, the command SHALL include all resource entries regardless of vote count.

#### Scenario: Zero-vote entries shown with --all
- **WHEN** "nodes" has 0 votes and `--all` is passed
- **THEN** "nodes" SHALL appear in the table output

### Requirement: Missing cache file handled gracefully
If `.ks.resources.json` does not exist, the command SHALL print an informative message and exit without error.

#### Scenario: No cache file
- **WHEN** `.ks.resources.json` does not exist in the configured directory
- **THEN** the command SHALL print a message indicating no resource data is found and exit with code 0
