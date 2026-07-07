## ADDED Requirements

### Requirement: Save scans all tmux panes
`ks save` SHALL enumerate all panes across all tmux sessions using `tmux list-panes -a`, collecting session name, window index, pane index, pane ID, pane PID, and current command for each pane.

#### Scenario: Panes enumerated successfully
- **WHEN** `ks save` is invoked inside an active tmux session
- **THEN** all panes across all sessions are listed with their positional indices and PIDs

#### Scenario: No tmux session
- **WHEN** `ks save` is invoked outside of tmux
- **THEN** the command exits cleanly with no output and no error

### Requirement: Save extracts KUBECONFIG from pane environment
For each enumerated pane, `ks save` SHALL read `/proc/<pane_pid>/environ` to extract the `KUBECONFIG` value. Panes with no `KUBECONFIG` in their process environment SHALL be silently skipped.

#### Scenario: Pane has KUBECONFIG set
- **WHEN** a pane's process environment contains `KUBECONFIG=/home/user/.kube/prod.yaml`
- **THEN** that path is recorded in the sidecar for that pane's positional index

#### Scenario: Pane has no KUBECONFIG
- **WHEN** a pane's process environment does not contain `KUBECONFIG`
- **THEN** the pane is omitted from the sidecar with no error

### Requirement: Save extracts cached resource from `.ks.resources.json`
For each pane with a `KUBECONFIG`, `ks save` SHALL look up the pane's ID in the resource cache (`~/.kube/.ks.resources.json`) using the current `$TMUX` session and the pane's ID. If a cached resource is found, it SHALL be included in the sidecar entry.

#### Scenario: Pane has cached resource
- **WHEN** the resource cache contains an entry for the pane's `$TMUX` + pane ID
- **THEN** the resource name (e.g., `pods`) is recorded in the sidecar

#### Scenario: Pane has no cached resource
- **WHEN** the resource cache has no entry for the pane
- **THEN** the sidecar entry omits the resource field; only KUBECONFIG is recorded

### Requirement: Save detects ks verb from k9s process arguments
If the pane's current command is `k9s`, `ks save` SHALL locate the k9s process in the pane's process tree, read its cmdline via `/proc/<k9s_pid>/cmdline`, and set the verb to `resource_all` if `-A` is present, otherwise `resource`.

#### Scenario: k9s running with -A flag
- **WHEN** `pane_current_command` is `k9s` and k9s cmdline contains `-A`
- **THEN** sidecar entry records `verb: resource_all`

#### Scenario: k9s running without -A flag
- **WHEN** `pane_current_command` is `k9s` and k9s cmdline does not contain `-A`
- **THEN** sidecar entry records `verb: resource`

#### Scenario: Pane not running k9s
- **WHEN** `pane_current_command` is not `k9s`
- **THEN** no verb is recorded; restore will use `resource` as default if resource is present

### Requirement: Save writes positional sidecar file
`ks save` SHALL write `~/.ks/.ks.resurrect.json` (respecting `--dir` / `KS_DIR`) keyed by stable positional indices (`session`, `window_idx`, `pane_idx`). The file SHALL include a `version` field.

#### Scenario: Sidecar written with populated panes
- **WHEN** at least one pane has KUBECONFIG set
- **THEN** `~/.ks/.ks.resurrect.json` is written with all qualifying pane entries

#### Scenario: Sidecar overwrites previous file
- **WHEN** `ks save` is run and a previous sidecar exists
- **THEN** the previous file is replaced with the current snapshot

### Requirement: Save is a no-op on non-Linux platforms
On platforms without `/proc`, `ks save` SHALL exit cleanly with no error and no output.

#### Scenario: Non-Linux platform
- **WHEN** `ks save` is run on macOS or another non-Linux OS
- **THEN** the command exits with code 0 and writes nothing
