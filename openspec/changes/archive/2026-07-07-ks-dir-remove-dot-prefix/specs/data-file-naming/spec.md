## ADDED Requirements

### Requirement: Data files in DataDir use non-hidden names
Data files written by `ks` to the `DataDir` (`~/.ks`) SHALL use plain filenames without a leading dot so they are visible by default with `ls` and standard file browsers.

The canonical filenames SHALL be:
- `resources.json` — k8s resource cache and vote leaderboard
- `resurrect.json` — tmux-resurrect session snapshot

#### Scenario: Resource cache file is visible
- **WHEN** `ks resource_load` writes the resource cache to `DataDir`
- **THEN** the file is created at `<DataDir>/resources.json` (no leading dot)

#### Scenario: Resurrect snapshot file is visible
- **WHEN** `ks save` writes the resurrect snapshot to `DataDir`
- **THEN** the file is created at `<DataDir>/resurrect.json` (no leading dot)

#### Scenario: Resource cache can be read back
- **WHEN** any `ks` command calls `resources.LoadResources(DataDir)`
- **THEN** the file is read from `<DataDir>/resources.json`

#### Scenario: Resurrect snapshot can be read back
- **WHEN** `ks restore` or `ks save` calls `resurrect.Load(DataDir)`
- **THEN** the file is read from `<DataDir>/resurrect.json`
