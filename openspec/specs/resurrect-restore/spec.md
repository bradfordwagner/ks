## ADDED Requirements

### Requirement: Restore reads sidecar and matches panes by position
`ks restore` SHALL read `~/.ks/.ks.resurrect.json` (respecting `--dir` / `KS_DIR`), enumerate current tmux panes, and match sidecar entries to live panes using `session_name + window_index + pane_index`.

#### Scenario: Sidecar entry matches a live pane
- **WHEN** a sidecar entry's session/window_idx/pane_idx corresponds to an existing pane
- **THEN** the restore action is applied to that pane

#### Scenario: Sidecar entry has no matching pane
- **WHEN** a sidecar entry's positional index does not match any current pane
- **THEN** the entry is skipped silently

#### Scenario: No sidecar file exists
- **WHEN** `ks restore` is invoked and `~/.ks/.ks.resurrect.json` does not exist
- **THEN** the command exits cleanly with no error

### Requirement: Restore sends inline KUBECONFIG and ks verb into pane
For each matched pane that has a resource in the sidecar, `ks restore` SHALL send `KS_RESOURCE=<resource> KUBECONFIG=<path> ks <verb>` as an inline command (not exported) via `tmux send-keys`. The verb defaults to `resource` if not recorded in the sidecar.

#### Scenario: Pane has KUBECONFIG and resource
- **WHEN** the sidecar entry contains both `kubeconfig` and `resource`
- **THEN** restore sends `KS_RESOURCE=pods KUBECONFIG=/home/user/.kube/prod.yaml ks resource` (or `ks resource_all`) into the pane via send-keys

#### Scenario: Pane has KUBECONFIG but no resource
- **WHEN** the sidecar entry contains `kubeconfig` but no `resource`
- **THEN** restore sends only `export KUBECONFIG=<path>` into the pane, leaving resource selection to fzf on next `ks resource` invocation

### Requirement: Restore does not export KS_RESOURCE into the shell environment
`KS_RESOURCE` SHALL be set only as an inline prefix on the `ks` invocation, never via `export`. This ensures `clear_pane` and `clear_cache` continue to work correctly — after k9s exits, `KS_RESOURCE` is absent from the pane's shell environment.

#### Scenario: User runs clear_pane after restore
- **WHEN** `ks restore` has sent `KS_RESOURCE=pods ks resource` into a pane and the user subsequently runs `ks clear_pane`
- **THEN** the pane's cache entry is cleared and the next `ks resource` invocation shows fzf

#### Scenario: User runs clear_cache after restore
- **WHEN** `ks restore` has sent inline commands into multiple panes and the user runs `ks clear_cache`
- **THEN** all cache entries are cleared and subsequent `ks resource` invocations show fzf in all panes
