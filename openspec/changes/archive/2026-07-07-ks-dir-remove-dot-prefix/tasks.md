## 1. Rename Constants

- [x] 1.1 In `internal/resources/types.go`, change `CacheFile = ".ks.resources.json"` to `CacheFile = "resources.json"`
- [x] 1.2 In `internal/resurrect/types.go`, change `ResurrectFile = ".ks.resurrect.json"` to `ResurrectFile = "resurrect.json"`

## 2. Fix Hardcoded Literal in Tests

- [x] 2.1 In `internal/cmds/resource_leaderboard_test.go` line 55, replace the hardcoded string `".ks.resources.json"` with `resources.CacheFile`

## 3. Update Documentation

- [x] 3.1 Update `CLAUDE.md` — replace all references to `.ks.resources.json` with `resources.json` and `.ks.resurrect.json` with `resurrect.json`

## 4. Verify

- [x] 4.1 Run `go test ./...` to confirm all tests pass with the new filenames
- [x] 4.2 Run `go vet ./...` to confirm no vet issues
