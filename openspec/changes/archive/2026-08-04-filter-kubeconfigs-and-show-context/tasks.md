## 1. Current-context resolution helper

- [x] 1.1 Add `CurrentContext(kubeconfig string) string` to `internal/kube/client.go`: `clientcmd.LoadFromFile`, return `load.CurrentContext` if non-empty, else `filepath.Base(kubeconfig)`; never return an error to the caller

## 2. Picker prompt header

- [x] 2.1 Change `internal/choose.One(opts []string, header string)` and `internal/choose.Multi(opts []string, header string)` to pass `fzf.WithPrompt(fmt.Sprintf("[%s] > ", header))` into the existing `choose()` helper alongside current options
- [x] 2.2 Color the prompt Catppuccin Mocha blue (`#89b4fa`) via `fzf.WithStyles(fzf.WithStylePrompt(fzf.Style{ForegroundColor: "#89b4fa"}))`, added to both `One` and `Multi`'s option list

## 3. Wire header into every call site

- [x] 3.1 `cmd/ks/main.go`: call `flag_helper.Load(&standardArgs)` in `rootCmd.RunE` before the command picker loop; pass `kube.CurrentContext(standardArgs.Kubeconfig)` as the header to `choose.One`
- [x] 3.2 `internal/cmds/set_ns_cmd.go`: pass `kube.CurrentContext(a.Kubeconfig)` to `choose.One`
- [x] 3.3 `internal/cmds/kube_new_ns_cmd.go`: pass `kube.CurrentContext(a.Kubeconfig)` to `choose.One`
- [x] 3.4 `internal/cmds/resource_cmd.go`: add `header string` parameter to `resolveResourceType`, use it only on the fzf-fallback `choose.One` call; update `Resource(a args.Standard, all bool)` to call `resolveResourceType(&r, kube.CurrentContext(a.Kubeconfig))`
- [x] 3.5 `internal/cmds/tmux_multi_cmd.go`: pass `kube.CurrentContext(a.Kubeconfig)` to `choose.Multi`
- [x] 3.6 `internal/cmds/tmux_window_cmd.go`: pass `kube.CurrentContext(a.Kubeconfig)` to `choose.One`
- [x] 3.7 `internal/cmds/kube_cp_cmd.go`: pass `kube.CurrentContext(a.Kubeconfig)` to `choose.One`
- [x] 3.8 `internal/cmds/link_cmd.go`: pass `kube.CurrentContext(a.Kubeconfig)` to `choose.One`
- [x] 3.9 `internal/cmds/pipe_cmd.go`: pass `kube.CurrentContext(a.Kubeconfig)` to `choose.One`

## 4. Test updates

- [x] 4.1 Update `internal/cmds/resource_cmd_test.go`: add a header argument (e.g. `""`) to existing `resolveResourceType(r)` calls

## 5. Verification

- [x] 5.1 `go build ./...` passes (pre-existing `go vet` lostcancel warnings in `set_ns_cmd.go`/`kube_new_ns_cmd.go` predate this change)
- [x] 5.2 `go test ./...` passes
- [x] 5.3 Manual smoke test: confirmed `kube.CurrentContext` resolves `current-context` from a kubeconfig and falls back to the base filename when the file is missing
