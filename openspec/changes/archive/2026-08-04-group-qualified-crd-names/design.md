## Context

`ResourceLoad` (`internal/cmds/resource_load_cmd.go`) walks the discovery API's `[]*APIResourceList`, and for each `APIResource` stores `apiResource.Name` (deduped, sub-resources skipped). The group — available on each list's `GroupVersion` — is dropped today. The fzf wrapper (`internal/choose`) matches on the same string it displays (go-fzf's `itemFunc` is both search key and label; no hidden key), so any group text must be part of the stored name. See proposal.md for motivation.

## Goals / Non-Goals

**Goals:**
- CRD names carry their group so they are searchable and unambiguous.
- Built-in resources stay short (per the chosen scope).
- Votes survive the one-time name shape change.
- No changes to `choose`, `resource`, or `resource_multi` — they already pass the name straight through.

**Non-Goals:**
- Putting the API version in the name.
- Renaming core/standard resources.
- Any per-group configuration surface.
- Changing k9s.

## Decisions

### D1: "Standard group" heuristic (keep built-ins short)

Qualify a name only when its group is non-empty **and** not standard. A group is standard when it ends with `.k8s.io` **or** is one of `{apps, batch, autoscaling, extensions, policy}`. Everything else — vendor domains like `redhatcop.redhat.io`, `cert-manager.io`, `argoproj.io` — is qualified.

- **Why:** official Kubernetes groups are either these few short names or `*.k8s.io`; operator CRDs use vendor DNS domains. This captures the user's "only CRDs / non-standard groups" scope without maintaining an exhaustive allowlist of every k8s group.
- **Alternatives:** (a) explicit allowlist of every standard group — brittle and long; rejected. (b) qualify everything (`deployments.apps`) — rejected per the scope decision to keep built-ins short.
- **Trade-off:** a CRD served under a `*.k8s.io` group (rare) would stay short. Acceptable.

### D2: Qualify with group only, not version

Produce `plural.group` (`randomsecrets.redhatcop.redhat.io`), never `plural.version.group`. Matches kubectl's fully-qualified form and keeps one entry per resource across served versions (the existing dedup Set still collapses them).

### D3: Preserve votes across the rename

When building merged entries, look up the vote count by the new qualified name first, then fall back to the old bare name. A one-time carry-over so the leaderboard and vote-sort don't reset for CRDs on the first re-load.

### D4: Pure, testable helpers

Factor the naming into `isStandardGroup(group string) bool` and `qualifiedResourceName(name, group string) string` so the behavior is unit-tested without a live cluster. Parse the group from `APIResourceList.GroupVersion` via `schema.ParseGroupVersion`.

## Risks / Trade-offs

- **k9s must accept `-c <plural>.<group>`.** → k9s resolves the GVR plural.group form; verify on a live cluster during apply. If some k9s build rejects it, the fallback is to keep the qualified name for display/search/cache but strip to the bare plural only for the `-c` argument — isolated to `resource_cmd.go`.
- **`resources.json` CRD names change shape.** → Votes carried by D3; the per-pane cache keeps any old bare entries, which still relaunch fine (k9s accepts bare CRD names too) and re-cache on next selection. No migration script needed.
- **Rare CRD under a `*.k8s.io` group stays short.** → Minor; does not break anything, just isn't group-filterable.
