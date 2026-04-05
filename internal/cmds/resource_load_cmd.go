package cmds

import (
	"errors"
	"github.com/bradfordwagner/go-util/bwutil"
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/kube"
	"github.com/bradfordwagner/ks/internal/resources"
	"os"
	"sort"
	"strings"
)

func ResourceLoad(a args.Standard) (err error) {
	l := log.Log()

	client, err := kube.Client(a.Kubeconfig)
	if err != nil {
		return err
	}

	_, kuberesources, err := client.Discovery().ServerGroupsAndResources()
	if err != nil {
		return err
	}

	// Load existing file to preserve votes; ignore not-found.
	existing, loadErr := resources.LoadResources(a.Directory)
	if loadErr != nil && !errors.Is(loadErr, os.ErrNotExist) {
		// Non-fatal: log and proceed with empty votes.
		l.With("err", loadErr).Warn("could not load existing resources file; votes will be reset")
	}

	// Build votes map from existing entries.
	votes := make(map[string]int, len(existing.Names))
	for _, e := range existing.Names {
		votes[e.Name] = e.Votes
	}

	// Collect API resource names (deduplicated, no sub-resources).
	resourceNames := bwutil.NewSet[string]()
	for _, resource := range kuberesources {
		for _, apiResource := range resource.APIResources {
			if !strings.Contains(apiResource.Name, "/") {
				resourceNames.Add(apiResource.Name)
			}
		}
	}
	keys := resourceNames.Keyset()
	sort.Strings(keys)

	// Construct merged v2 entries (pruning names not returned by the API).
	entries := make([]resources.ResourceEntry, 0, len(keys))
	for _, name := range keys {
		l.With("resource", name).Info("found resource")
		entries = append(entries, resources.ResourceEntry{
			Name:  name,
			Votes: votes[name],
		})
	}

	r := resources.Resources{
		Version: resources.CurrentVersion,
		Cache:   existing.Cache,
		Names:   entries,
	}

	resourceFile, err := r.Write(a.Directory)
	if err != nil {
		return err
	}
	l.With("file", resourceFile).Info("wrote resources to file")

	return
}
