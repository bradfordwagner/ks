package cmds

import (
	"errors"

	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/bradfordwagner/ks/internal/kube"
	"github.com/bradfordwagner/ks/internal/resources"
)

// ResourceExists is the instant counterpart to ResourceLive: it filters the
// picker to resource types recorded as populated at the last resource_load probe
// (no live cluster calls). Selection behaves like Resource — one pick opens in
// the current pane, more open a tiled split pane each.
func ResourceExists(a args.Standard, all bool) (err error) {
	l := log.Log()

	r, err := resources.LoadResources(a.DataDir)
	if err != nil {
		return
	}

	names := populatedNames(&r)
	if len(names) == 0 {
		l.Info("no populated resources recorded; run `ks resources_load` to probe usage")
		return nil
	}

	selected, err := choose.Multi(names, kube.CurrentContext(a.Kubeconfig))
	if errors.Is(err, choose.ErrAbort) {
		return nil
	}
	if err != nil {
		return
	}
	return launchSelected(a, &r, selected, all)
}

// populatedNames returns the cached resource names known to have objects at the
// last resource_load probe (Probed && Populated), vote-sorted (most-used first).
func populatedNames(r *resources.Resources) []string {
	populated := make(map[string]bool)
	for _, e := range r.Names {
		if e.Probed && e.Populated {
			populated[e.Name] = true
		}
	}
	out := make([]string, 0, len(populated))
	for _, name := range r.SortedNames() {
		if populated[name] {
			out = append(out, name)
		}
	}
	return out
}
