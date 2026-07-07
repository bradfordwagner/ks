package cmds

import (
	"errors"
	"os"

	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/bradfordwagner/ks/internal/k9s"
	"github.com/bradfordwagner/ks/internal/resources"
	"github.com/bradfordwagner/ks/internal/resurrect"
	"github.com/koki-develop/go-fzf"
)

// Resource opens k9s with selected resource view
func Resource(a args.Standard, all bool) (err error) {
	r, err := resources.LoadResources(a.DataDir)
	if err != nil {
		return
	}

	resourceType, err := resolveResourceType(&r)
	if err != nil || resourceType == "" {
		return
	}

	verb := "resource"
	if all {
		verb = "resource_all"
	}
	r.VoteFor(resourceType)
	go r.Write(a.DataDir)
	upsertResurrectPane(a.DataDir, resourceType, verb)

	k9sArgs := []string{"-c", resourceType}
	if all {
		k9sArgs = append(k9sArgs, "-A")
	}
	k9s.Run(k9sArgs...)
	return
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

// resolveResourceType returns the resource to display in k9s.
// Priority: KS_RESOURCE env var → pane cache → fzf selection.
// Returns ("", nil) if the user aborted fzf.
func resolveResourceType(r *resources.Resources) (resourceType string, err error) {
	if override := os.Getenv("KS_RESOURCE"); override != "" {
		r.Upsert(override)
		return override, nil
	}
	resourceType = r.Get()
	if resourceType == "" {
		resourceType, err = choose.One(r.SortedNames())
		if errors.Is(err, fzf.ErrAbort) {
			return "", nil
		}
		if err != nil {
			return "", err
		}
		r.Upsert(resourceType)
	}
	return
}
