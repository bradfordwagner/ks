## Why

Files in `~/.ks` use a leading dot (e.g., `.ks.resources.json`, `.ks.resurrect.json`), making them hidden by default in most file browsers and `ls` output without `-a`. Since `~/.ks` is already a dedicated directory, the dot prefix adds no value and makes manual inspection and debugging harder.

## What Changes

- Rename `.ks.resources.json` → `resources.json` in the `internal/resources` package
- Rename `.ks.resurrect.json` → `resurrect.json` in the `internal/resurrect` package
- Update any inline documentation (comments, README) referencing the old names

## Capabilities

### New Capabilities

- `data-file-naming`: Data files stored under `~/.ks` use plain (non-hidden) filenames for visibility

### Modified Capabilities

<!-- No spec-level behavioral requirements are changing — this is a filename-only implementation change -->

## Impact

- `internal/resources/types.go`: `CacheFile` constant
- `internal/resurrect/types.go`: `ResurrectFile` constant
- CLAUDE.md documentation references to old file names
- **BREAKING**: Existing `~/.ks/.ks.resources.json` and `~/.ks/.ks.resurrect.json` files on disk will not be auto-migrated; users must rename them manually or re-run `ks resource_load` and `ks save`
