## Purpose

Defines how the resource picker turns a selection of resource types into k9s views, including multi-selection that tiles one k9s pane per additional pick within the current cluster.

## ADDED Requirements

### Requirement: Multi-select launches a pane per selection

When the picker is shown, `resource` (and `resource_all`) SHALL allow selecting multiple resource types. The first selected type SHALL open k9s in the current pane; each additional selected type SHALL open in its own tiled tmux split pane. A single selection therefore behaves exactly as the original single-resource flow, with no splits.

#### Scenario: A single selection opens in the current pane
- **WHEN** the user selects exactly one resource type
- **THEN** k9s opens on that type in the current pane and no split panes are created

#### Scenario: Multiple selections tile a pane each
- **WHEN** the user selects three resource types
- **THEN** the first opens k9s in the current pane and the other two each open in a tiled split pane running k9s on their type

### Requirement: Split panes share the invoking cluster

Each split pane opened for an additional selection SHALL inherit the invoking pane's `KUBECONFIG`, so every pane targets the same cluster.

#### Scenario: Split panes follow the current kubeconfig
- **WHEN** the invoking shell has `KUBECONFIG` set and the user selects two resource types
- **THEN** the split pane opened for the second type has that same `KUBECONFIG` in its environment

### Requirement: Abort is a no-op

An aborted picker (Esc / Ctrl-C) SHALL open no panes and return no error.

#### Scenario: User aborts the picker
- **WHEN** the user dismisses the picker without confirming a selection
- **THEN** no panes are opened and the command exits without error

### Requirement: All-namespaces via the resource_all verb

`resource_all` SHALL behave like `resource` but watch every selected type across all namespaces; split panes opened by `resource_all` SHALL themselves run the all-namespaces verb.

#### Scenario: resource_all watches all namespaces for every selection
- **WHEN** the user runs `resource_all` and selects two types
- **THEN** the current pane and the split pane both watch their type across all namespaces

### Requirement: A resolved single resource skips the picker

When a single resource is already determined for the pane — via the `KS_RESOURCE` override or the pane cache (see the ks-resource-env-override capability) — the picker SHALL NOT be shown and k9s SHALL open that resource in the current pane. This keeps split panes and resurrect restore (which set `KS_RESOURCE`) launching a single view.

#### Scenario: KS_RESOURCE bypasses the multi-select picker
- **WHEN** `KS_RESOURCE=pods ks resource` is invoked
- **THEN** k9s opens on `pods` in the current pane without showing the picker

#### Scenario: A cached pane reopens its resource
- **WHEN** `ks resource` runs in a pane that already has a cached resource
- **THEN** k9s reopens that resource in the current pane without showing the picker
