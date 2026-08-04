## 1. Multi-select resource selection

- [x] 1.1 Add `tmux.SplitAndSend(kubeconfig, keys)` (shared tiled-split helper) alongside `Split`.
- [x] 1.2 Rewrite `Resource` to resolve a single resource from `KS_RESOURCE`/cache first (`singleResource`), else multi-select via `choose.Multi`.
- [x] 1.3 Add `launchSelected` (first pick → current pane, rest → tiled splits) and `runSingle` (cache/vote/resurrect/k9s), shared by all pickers.
- [x] 1.4 Remove the `resource_multi` command (cmd + cmd file + test) and drop it from the root menu.

## 2. Usage probe

- [x] 2.1 Add `internal/cmds/discovery.go`: `listableGVRs` (name→GVR from discovery) and `probePopulated` (parallel `Limit:1`, all namespaces, bounded concurrency).
- [x] 2.2 Add `Populated`/`Probed` to `ResourceEntry`; bump schema to v3 (additive load).
- [x] 2.3 In `resource_load`, build the GVR map inline, probe existence (bounded by `--timeout`, non-fatal), and store the flags per entry.

## 3. resource_exists (cached) and resource_live (live)

- [x] 3.1 Add `resource_exists` (`populatedNames` filters cached `Probed && Populated`, vote-sorted; empty-state guidance) + cmd wiring + `-A`.
- [x] 3.2 Add `resource_live` (batch parallel probe → `prioritizeCandidates` ordering → `choose.Multi` → `launchSelected`) + cmd wiring + `-A`.
- [x] 3.3 Register both in the root menu.
- [x] 3.4 Stream `resource_live` results live via `choose.MultiStream` + `probePopulatedStream` (priority-dispatched channel), rendered in the tmux popup (see the fzf-binary-picker change) so incremental updates render cleanly. Earlier full-screen/go-fzf streaming attempts corrupted (design.md D4).

## 4. Tests & docs

- [x] 4.1 Unit tests: `singleResource`, `resourceVerb`/`resourcePaneCommand`, `prioritizeCandidates`, `populatedNames`.
- [x] 4.2 Update `CLAUDE.md` (resource multi-select, resource_live, resource_exists, resource_load probe, v3 cache; remove resource_multi).
- [x] 4.3 `go vet` / `go test ./...` / `go install ./cmd/ks`.

## 5. Verification (pending)

- [ ] 5.1 Manually verify against a live cluster: `ks resource` multi-select tiles panes; `ks resource_exists` shows cached-populated; `ks resource_live` live-scans and prioritizes voted+previously-populated; `-A` variants watch all namespaces.
