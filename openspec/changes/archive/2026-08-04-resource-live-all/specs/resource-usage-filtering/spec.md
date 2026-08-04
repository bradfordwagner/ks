## ADDED Requirements

### Requirement: All-namespaces via the resource_live_all verb

`resource_live_all` SHALL behave like `resource_live` — live-scanning the cluster and offering only currently-populated types, streamed into the picker as found — but SHALL watch every selected type across all namespaces, equivalent to running `resource_live -A`.

#### Scenario: resource_live_all watches all namespaces for the selection
- **WHEN** the user runs `resource_live_all` and selects a type
- **THEN** the live scan behaves as it does for `resource_live`, and the selected type opens k9s watching that type across all namespaces

#### Scenario: resource_live_all is reachable from the interactive menu
- **WHEN** the user runs bare `ks` and picks `resource_live_all` from the command menu
- **THEN** it runs the same all-namespaces live scan as `ks resource_live -A`, without requiring any flag to be passed
