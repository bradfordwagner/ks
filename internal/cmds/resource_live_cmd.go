package cmds

import (
	"context"
	"errors"
	"sort"

	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/bradfordwagner/ks/internal/kube"
	"github.com/bradfordwagner/ks/internal/resources"
)

// ResourceLive live-scans the cluster for resource types that currently have
// objects and lets the user multi-select from only those. Results stream into the
// picker (fzf reads its input incrementally) as the parallel scan finds them;
// probes are dispatched so types that were populated at the last resource_load and
// are in the leaderboard (voted) are checked — and thus tend to appear — first.
// Inside tmux the picker renders in a popup so incremental updates render cleanly
// (see internal/choose). Selection behaves like Resource: one pick opens in the
// current pane, more open a tiled split pane each.
func ResourceLive(a args.Standard, all bool) (err error) {
	r, err := resources.LoadResources(a.DataDir)
	if err != nil {
		return
	}

	client, err := kube.Client(a.Kubeconfig)
	if err != nil {
		return
	}
	gvrByName, err := listableGVRs(client.Discovery())
	if err != nil {
		return
	}
	dyn, err := kube.Dynamic(a.Kubeconfig)
	if err != nil {
		return
	}

	candidates := make([]string, 0, len(gvrByName))
	for name := range gvrByName {
		candidates = append(candidates, name)
	}
	order := prioritizeCandidates(r.Names, candidates)

	ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()

	// Stream populated types into the picker as the scan finds them.
	names := make(chan string)
	go func() {
		probePopulatedStream(ctx, dyn, gvrByName, order, func(name string) {
			select {
			case names <- name:
			case <-ctx.Done():
			}
		})
		close(names)
	}()

	selected, err := choose.MultiStream(names, kube.CurrentContext(a.Kubeconfig))
	cancel() // stop stray probes before launching k9s
	if errors.Is(err, choose.ErrAbort) {
		return nil
	}
	if err != nil {
		return
	}
	return launchSelected(a, &r, selected, all)
}

// prioritizeCandidates orders resource names so that types which were populated at
// the last resource_load AND have votes come first (votes descending, then name),
// then the rest (same sort). Names absent from the cache sort as zero-vote,
// never-populated. Used to order the resource_live picker (most relevant first).
func prioritizeCandidates(entries []resources.ResourceEntry, candidates []string) []string {
	byName := make(map[string]resources.ResourceEntry, len(entries))
	for _, e := range entries {
		byName[e.Name] = e
	}

	type info struct {
		name  string
		votes int
	}
	var priority, rest []info
	for _, name := range candidates {
		e := byName[name] // zero value when the candidate isn't in the cache
		it := info{name: name, votes: e.Votes}
		if e.Votes > 0 && e.Populated {
			priority = append(priority, it)
		} else {
			rest = append(rest, it)
		}
	}

	byVotesThenName := func(s []info) {
		sort.Slice(s, func(i, j int) bool {
			if s[i].votes != s[j].votes {
				return s[i].votes > s[j].votes
			}
			return s[i].name < s[j].name
		})
	}
	byVotesThenName(priority)
	byVotesThenName(rest)

	out := make([]string, 0, len(candidates))
	for _, it := range priority {
		out = append(out, it.name)
	}
	for _, it := range rest {
		out = append(out, it.name)
	}
	return out
}
