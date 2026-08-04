## 1. Naming logic

- [x] 1.1 Add pure helpers in `internal/cmds` — `isStandardGroup(group string) bool` and `qualifiedResourceName(name, group string) string` — implementing D1 (standard = empty / `*.k8s.io` / `{apps,batch,autoscaling,extensions,policy}`) and D2 (group only, no version).
- [x] 1.2 In `ResourceLoad`, parse each list's group via `schema.ParseGroupVersion(resource.GroupVersion)` and build the stored name with `qualifiedResourceName`; keep the sub-resource skip (`/`) and dedup Set.
- [x] 1.3 Preserve votes (D3): look up the vote count by the qualified name, then fall back to the bare name.

## 2. Tests

- [x] 2.1 Add a Ginkgo spec covering `isStandardGroup` (core `""`, `apps`, `batch`, `foo.k8s.io` → true; `redhatcop.redhat.io`, `cert-manager.io` → false) and `qualifiedResourceName` (core/standard → bare, `redhatcop.redhat.io` → `randomsecrets.redhatcop.redhat.io`).

## 3. Docs & verification

- [x] 3.1 Update `CLAUDE.md` (`resource_load` / `resource` reference and the Resource Cache section) to note CRD names are group-qualified and filterable by typing the group.
- [x] 3.2 Run `go vet` on changed packages, `go test ./...`, and `go install ./cmd/ks`.
- [ ] 3.3 Manually verify against a cluster with vault-config-operator: `ks resource_load`, then `ks resource` → typing `redhatcop` filters to the VCO CRDs; selecting one launches k9s at it; `ks resource_multi` + `redhatcop` tiles them.
