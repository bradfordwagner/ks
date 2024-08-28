# ks

- [y] ks pipe - KUBECONFIG=$(ks pipe)
- [y] ks tmux_multi
- [y] ks tmux_window
  - [y] remove kcc hack
- [y] ks link
- [y] resources load
- existing commands
  - [y] resource
  - [x] ns_resource
  - [y] select_ns
  - [y] new_ctx_ns_cp (kube_new_ns)
  - [y] select_primary_kube_ctx
  - [y] select_kube_ctx_cp (kube_cp)
  - [y] multi_kube_ctx
  - [y] kube_new_window
  - [y] allns_resource (all_resource)
- [y] jump off to all cmds
  - [y] allow completion ordering for default


## Local
```bash
go install ./cmd/ks
```
