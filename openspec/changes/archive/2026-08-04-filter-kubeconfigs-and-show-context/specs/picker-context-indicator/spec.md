## ADDED Requirements

### Requirement: Every fzf picker displays the active kubeconfig context
Every interactive fzf picker launched by `ks` (root command menu, `set_ns`, `resource`/`resource_all`, `kube_new_ns`, `tmux_multi`, `tmux_window`, `kube_cp`, `link`, `pipe`) SHALL display the current-context name of the active kubeconfig (`a.Kubeconfig`) as a prefix on the picker's input prompt, colored Catppuccin Mocha blue (`#89b4fa`), so the active cluster is visible without leaving the picker.

#### Scenario: Kubeconfig has a resolvable current-context
- **WHEN** a picker is launched and the active kubeconfig file has a non-empty `current-context` field
- **THEN** the picker's prompt is rendered as `[<current-context>] > ` with the prompt text colored Catppuccin Mocha blue (`#89b4fa`)

#### Scenario: Kubeconfig is missing or unreadable
- **WHEN** a picker is launched and the active kubeconfig path cannot be loaded or parsed
- **THEN** the picker's prompt falls back to `[<base filename of the kubeconfig path>] > `, still colored Catppuccin Mocha blue, instead of erroring or omitting the indicator

### Requirement: `internal/kube.CurrentContext` resolves a display label for a kubeconfig
`internal/kube` SHALL expose a function that, given a kubeconfig path, returns the kubeconfig's `current-context` name, or the base filename if the context cannot be determined.

#### Scenario: Valid kubeconfig with a current-context set
- **WHEN** `CurrentContext` is called with a path to a valid kubeconfig file whose `current-context` field is `prod-cluster`
- **THEN** it returns `prod-cluster`

#### Scenario: Kubeconfig file does not exist
- **WHEN** `CurrentContext` is called with a path that does not exist on disk
- **THEN** it returns the base filename of the given path (no error is surfaced to the caller)
