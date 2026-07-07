## ADDED Requirements

### Requirement: KS_RESOURCE env var overrides fzf and cache lookup
If the `KS_RESOURCE` environment variable is set when `ks resource` or `ks resource_all` is invoked, the value SHALL be used as the resource type directly, bypassing both the cache lookup (`r.Get()`) and fzf selection.

#### Scenario: KS_RESOURCE is set
- **WHEN** `KS_RESOURCE=pods ks resource` is invoked
- **THEN** `pods` is used as the resource type without consulting the cache or launching fzf

#### Scenario: KS_RESOURCE is not set
- **WHEN** `ks resource` is invoked without `KS_RESOURCE`
- **THEN** behavior is unchanged: cache lookup first, fzf if cache miss

### Requirement: KS_RESOURCE invocation populates the cache
When `KS_RESOURCE` is used, `Resource()` SHALL call `r.Upsert(resourceType)` and `r.Write()` so the cache is populated for subsequent invocations under the pane's new `$TMUX_PANE` ID.

#### Scenario: Cache populated after KS_RESOURCE invocation
- **WHEN** `KS_RESOURCE=pods ks resource` runs and k9s exits
- **THEN** subsequent `ks resource` invocations in the same pane find a cache hit for `pods` without showing fzf

#### Scenario: Vote incremented for KS_RESOURCE resource
- **WHEN** `KS_RESOURCE=pods ks resource` runs
- **THEN** the vote count for `pods` in `.ks.resources.json` is incremented by one

### Requirement: KS_RESOURCE does not persist in the shell environment
`KS_RESOURCE` SHALL only be honored when set as an inline env var prefix on the `ks` invocation. The implementation SHALL NOT export or persist `KS_RESOURCE` itself. This ensures that after k9s exits, the pane's shell has no `KS_RESOURCE` set, and `clear_pane` / `clear_cache` restore normal fzf behavior.

#### Scenario: KS_RESOURCE absent after invocation
- **WHEN** the shell invokes `KS_RESOURCE=pods ks resource` and k9s exits
- **THEN** `echo $KS_RESOURCE` in the pane returns empty

#### Scenario: clear_pane works after KS_RESOURCE invocation
- **WHEN** `KS_RESOURCE=pods ks resource` has been used and the user runs `ks clear_pane`
- **THEN** the next `ks resource` invocation shows fzf
