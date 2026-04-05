## 1. Schema & Types

- [x] 1.1 Add `Version int` and `ResourceEntry{Name, Votes}` to `internal/resources/types.go`; change `Names` field type from `[]string` to `[]ResourceEntry`
- [x] 1.2 Implement two-pass v1→v2 migration in `LoadResources`: decode `names` as `json.RawMessage`, attempt `[]ResourceEntry`, fall back to `[]string` with votes initialised to 0; set `Version = 2`
- [x] 1.3 Add `VoteFor(name string)` method on `Resources` that increments the matching entry's `votes` by 1
- [x] 1.4 Add `SortedNames() []string` method that returns resource names sorted by votes descending then alphabetically ascending

## 2. Resource Load Command

- [x] 2.1 In `ResourceLoad`, attempt to load existing `.ks.resources.json` before querying the API (ignore not-found error)
- [x] 2.2 Build a `map[string]int` of existing votes from the loaded file
- [x] 2.3 After collecting API resource names, construct `[]ResourceEntry` using the votes map (0 for new names); entries not returned by the API are pruned
- [x] 2.4 Write the merged v2 `Resources` struct back to disk

## 3. Resource Selection Commands

- [x] 3.1 In `Resource()`, call `r.VoteFor(resourceType)` after the user selects a resource type (both when chosen via fzf and when retrieved from pane cache)
- [x] 3.2 Pass `r.SortedNames()` to `choose.One()` instead of `r.Names` directly
- [x] 3.3 Ensure the async `r.Write(a.Directory)` call is triggered after both `Upsert` and `VoteFor`

## 4. Tests

- [x] 4.1 Unit test v1→v2 migration: load a v1 JSON fixture, assert all entries have `votes: 0` and `Version == 2`
- [x] 4.2 Unit test `VoteFor`: verify vote count increments and unknown name is a no-op
- [x] 4.3 Unit test `SortedNames`: verify descending-vote then ascending-alpha ordering
- [x] 4.4 Unit test `ResourceLoad` merge: existing votes preserved, new resources added at 0, removed resources pruned
