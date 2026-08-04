## 1. Command

- [x] 1.1 Create `cmd/ks/resource_live_all.go` mirroring `cmd/ks/resource_all.go`: register standard flags, define `resourceLiveAllCmd` with `Use: "resource_live_all"`, and `RunE` calling `cmds.ResourceLive(standardArgs, true)`.

## 2. Registration

- [x] 2.1 Add `resourceLiveAllCmd` to the `commands` slice in `cmd/ks/main.go` (next to `resourceLiveCmd`) so it appears in the interactive menu and as a top-level cobra command.

## 3. Documentation

- [x] 3.1 Add a `resource_live_all` row to the Command Reference table in `CLAUDE.md`, describing it as the all-namespaces counterpart to `resource_live` (equivalent to `resource_live -A`), reachable from the interactive menu.

## 4. Verification

- [x] 4.1 Run `go build ./...` and `go vet ./...`. (build passes cleanly; the two `go vet` context-leak findings in `internal/cmds/kube_new_ns_cmd.go` and `internal/cmds/set_ns_cmd.go` are pre-existing on `main`, unrelated to this change)
- [x] 4.2 Ran `go install ./cmd/ks`; confirmed `resource_live_all` is listed in `ks --help` (and therefore the bare-`ks` interactive menu, which is built from the same `commands` slice) and `ks resource_live_all --help` shows the standard flags. Full interactive run (fzf popup + k9s TUI against the live k3d-dev cluster) isn't drivable from this non-TTY shell; traced the code path instead — `resource_live_all` calls `cmds.ResourceLive(standardArgs, true)`, the exact same unmodified function and `all=true` value that `resource_live -A` already exercises, so its all-namespaces launch behavior is proven by the existing, already-shipped `resource_live -A` path.
