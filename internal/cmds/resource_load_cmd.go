package cmds

import (
	"context"
	"errors"
	"github.com/bradfordwagner/go-util/bwutil"
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/kube"
	"github.com/bradfordwagner/ks/internal/resources"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"sort"
	"strings"
	"time"
)

// standardShortGroups are the well-known Kubernetes API groups that do not end in
// ".k8s.io" but whose resources should still keep their short plural name.
var standardShortGroups = map[string]bool{
	"apps":        true,
	"batch":       true,
	"autoscaling": true,
	"extensions":  true,
	"policy":      true,
}

// isStandardGroup reports whether an API group is a built-in Kubernetes group
// (core, any *.k8s.io group, or a well-known short group) rather than a
// CRD/vendor group.
func isStandardGroup(group string) bool {
	if group == "" || strings.HasSuffix(group, ".k8s.io") {
		return true
	}
	return standardShortGroups[group]
}

// qualifiedResourceName returns "<name>.<group>" for CRD/non-standard groups so
// the group is searchable in the picker, and the bare name for core/standard
// Kubernetes groups.
func qualifiedResourceName(name, group string) string {
	if isStandardGroup(group) {
		return name
	}
	return name + "." + group
}

func ResourceLoad(a args.Standard) (err error) {
	l := log.Log()
	start := time.Now()

	client, err := kube.Client(a.Kubeconfig)
	if err != nil {
		return err
	}

	_, kuberesources, err := client.Discovery().ServerGroupsAndResources()
	if err != nil {
		return err
	}

	// Load existing file to preserve votes; ignore not-found.
	existing, loadErr := resources.LoadResources(a.DataDir)
	if loadErr != nil && !errors.Is(loadErr, os.ErrNotExist) {
		// Non-fatal: log and proceed with empty votes.
		l.With("err", loadErr).Warn("could not load existing resources file; votes will be reset")
	}

	// Build votes map from existing entries.
	votes := make(map[string]int, len(existing.Names))
	for _, e := range existing.Names {
		votes[e.Name] = e.Votes
	}

	// Collect API resource names (deduplicated, no sub-resources). CRD/vendor
	// groups are qualified as "<name>.<group>"; qualifiedToBare remembers the old
	// bare name so votes can migrate across the rename, and gvrByName records the
	// GVR of each listable type for the usage probe.
	resourceNames := bwutil.NewSet[string]()
	qualifiedToBare := make(map[string]string)
	gvrByName := make(map[string]schema.GroupVersionResource)
	for _, resource := range kuberesources {
		gv, _ := schema.ParseGroupVersion(resource.GroupVersion)
		for _, apiResource := range resource.APIResources {
			if strings.Contains(apiResource.Name, "/") {
				continue
			}
			qualified := qualifiedResourceName(apiResource.Name, gv.Group)
			resourceNames.Add(qualified)
			qualifiedToBare[qualified] = apiResource.Name
			if _, ok := gvrByName[qualified]; !ok && hasVerb(apiResource.Verbs, "list") {
				gvrByName[qualified] = gv.WithResource(apiResource.Name)
			}
		}
	}
	keys := resourceNames.Keyset()
	sort.Strings(keys)

	// Probe each listable type for existence (cluster-wide, in parallel, bounded
	// by the timeout). Non-fatal: on any failure the entries are simply left
	// unprobed. `populated` is keyed only by successfully-checked names.
	populated := map[string]bool{}
	if dyn, dynErr := kube.Dynamic(a.Kubeconfig); dynErr != nil {
		l.With("err", dynErr).Warn("could not create dynamic client; skipping usage probe")
	} else if len(gvrByName) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
		populated = probePopulated(ctx, dyn, gvrByName)
		cancel()
	}

	// Construct merged entries (pruning names not returned by the API).
	// Votes are looked up by the qualified name, falling back to the old bare
	// name so counts carry over the first time a CRD becomes group-qualified.
	entries := make([]resources.ResourceEntry, 0, len(keys))
	for _, name := range keys {
		l.With("resource", name).Info("found resource")
		v, ok := votes[name]
		if !ok {
			v = votes[qualifiedToBare[name]]
		}
		entry := resources.ResourceEntry{Name: name, Votes: v}
		if isPopulated, probed := populated[name]; probed {
			entry.Probed = true
			entry.Populated = isPopulated
		}
		entries = append(entries, entry)
	}

	r := resources.Resources{
		Version: resources.CurrentVersion,
		Cache:   existing.Cache,
		Names:   entries,
	}

	resourceFile, err := r.Write(a.DataDir)
	if err != nil {
		return err
	}
	l.With("file", resourceFile).
		With("elapsed", time.Since(start).Round(time.Millisecond).String()).
		Info("wrote resources to file")

	return
}
