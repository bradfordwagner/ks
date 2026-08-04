package cmds

import (
	"errors"
	"fmt"
	"os"

	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/bradfordwagner/ks/internal/k9s"
	"github.com/bradfordwagner/ks/internal/kube"
	"github.com/bradfordwagner/ks/internal/resources"
	"github.com/bradfordwagner/ks/internal/resurrect"
	"github.com/bradfordwagner/ks/internal/tmux"
)

// Resource opens k9s at a selected resource view.
//
// Fast paths (single resource, current pane): a KS_RESOURCE override (set by
// split panes and restore) or an existing pane cache entry launches k9s directly.
// Otherwise the user multi-selects from the picker: a single pick opens in the
// current pane, while multiple picks open the first in the current pane and one
// tiled split pane per remaining selection.
func Resource(a args.Standard, all bool) (err error) {
	r, err := resources.LoadResources(a.DataDir)
	if err != nil {
		return
	}

	if single, ok := singleResource(&r); ok {
		return runSingle(a, &r, single, all)
	}

	selected, err := choose.Multi(r.SortedNames(), kube.CurrentContext(a.Kubeconfig))
	if errors.Is(err, choose.ErrAbort) {
		return nil
	}
	if err != nil {
		return
	}
	return launchSelected(a, &r, selected, all)
}

// singleResource returns a resource resolved without prompting — from the
// KS_RESOURCE env override (cached to the pane) or the pane cache. ok is false
// when neither is set and the caller should prompt.
func singleResource(r *resources.Resources) (resourceType string, ok bool) {
	if override := os.Getenv("KS_RESOURCE"); override != "" {
		r.Upsert(override)
		return override, true
	}
	if cached := r.Get(); cached != "" {
		return cached, true
	}
	return "", false
}

// launchSelected opens k9s for a multi-selection: the first resource runs in the
// current pane, each remaining resource opens a tiled split pane that re-invokes
// ks for itself. A single selection therefore behaves exactly like the original
// single-resource flow (no splits).
func launchSelected(a args.Standard, r *resources.Resources, selected []string, all bool) error {
	if len(selected) == 0 {
		return nil
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	verb := resourceVerb(all)
	for _, resource := range selected[1:] {
		if err := tmux.SplitAndSend(kubeconfig, resourcePaneCommand(resource, verb)); err != nil {
			return err
		}
	}
	return runSingle(a, r, selected[0], all)
}

// runSingle caches/votes for the resource, records it for resurrect, and launches
// k9s in the current pane (blocking).
func runSingle(a args.Standard, r *resources.Resources, resourceType string, all bool) error {
	r.Upsert(resourceType)
	r.VoteFor(resourceType)
	go r.Write(a.DataDir)
	upsertResurrectPane(a.DataDir, resourceType, resourceVerb(all))

	k9sArgs := []string{"-c", resourceType}
	if all {
		k9sArgs = append(k9sArgs, "-A")
	}
	k9s.Run(k9sArgs...)
	return nil
}

// resourceVerb maps the all-namespaces toggle to the ks verb a split pane runs.
func resourceVerb(all bool) string {
	if all {
		return "resource_all"
	}
	return "resource"
}

// resourcePaneCommand builds the inline command a spawned pane runs to launch k9s
// for the resource. KS_RESOURCE is set inline (not exported) so the pane's own ks
// invocation caches and votes for it, and clear_pane / clear_cache keep working.
func resourcePaneCommand(resource, verb string) string {
	return fmt.Sprintf("KS_RESOURCE=%s ks %s", resource, verb)
}

// upsertResurrectPane writes the current pane's resource+verb into the resurrect sidecar.
func upsertResurrectPane(dataDir, resource, verb string) {
	pane, err := resurrect.CurrentPane()
	if err != nil || pane.Session == "" {
		return
	}
	_ = resurrect.Upsert(dataDir, resurrect.ResurrectPane{
		Session:    pane.Session,
		WindowIdx:  pane.WindowIdx,
		PaneIdx:    pane.PaneIdx,
		Kubeconfig: os.Getenv("KUBECONFIG"),
		Resource:   resource,
		Verb:       verb,
	})
}
