## Purpose

Defines how ks names Kubernetes resource types discovered from the cluster so the picker is searchable by API group and unambiguous across groups.

## ADDED Requirements

### Requirement: CRD resources are qualified by API group

When loading resource types from the cluster, ks SHALL name a resource from a custom / non-standard API group as `<plural>.<group>`. Resources in the core group and in standard Kubernetes groups SHALL keep their short plural name. A group is standard when it is empty (core), ends with `.k8s.io`, or is one of the well-known short groups (`apps`, `batch`, `autoscaling`, `extensions`, `policy`).

#### Scenario: Vault Config Operator CRD is group-qualified
- **WHEN** `resource_load` discovers `randomsecrets` in group `redhatcop.redhat.io`
- **THEN** the stored and displayed name is `randomsecrets.redhatcop.redhat.io`

#### Scenario: Core resource stays short
- **WHEN** `resource_load` discovers `pods` in the core (empty) group
- **THEN** the stored and displayed name is `pods`

#### Scenario: Standard Kubernetes group stays short
- **WHEN** `resource_load` discovers `deployments` in group `apps`, or any resource in a `*.k8s.io` group
- **THEN** the stored and displayed name keeps its short plural (e.g. `deployments`)

### Requirement: Group filtering in the picker

Because the qualified name contains the group, the picker SHALL allow narrowing to a group's resources by typing a substring of the group name.

#### Scenario: Typing the group filters to its CRDs
- **WHEN** the user types `redhatcop` in the resource picker
- **THEN** the visible options are limited to resources whose qualified name contains `redhatcop.redhat.io`

### Requirement: Qualified name used consistently downstream

The name produced at load time SHALL be the same value cached per pane, voted on, and passed to k9s, so that selection, caching, and launching all use one identifier.

#### Scenario: Selecting a qualified resource launches k9s at it
- **WHEN** the user selects `randomsecrets.redhatcop.redhat.io` from the picker
- **THEN** k9s is launched focused on that resource and that same qualified name is recorded as the pane's cached resource

### Requirement: Votes preserved across qualification

`resource_load` SHALL preserve a resource's existing selection votes when its stored name changes from the bare plural to the group-qualified form.

#### Scenario: Votes carry over on re-load
- **WHEN** a resource previously stored as `randomsecrets` with N votes is re-loaded and becomes `randomsecrets.redhatcop.redhat.io`
- **THEN** the qualified entry retains N votes
