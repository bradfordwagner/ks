package cmds

import (
	"context"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

// probeConcurrency bounds how many existence checks run at once so a large
// cluster's resource catalog doesn't flood the API server.
const probeConcurrency = 16

// hasVerb reports whether a resource's verb list contains v.
func hasVerb(verbs metav1.Verbs, v string) bool {
	for _, got := range verbs {
		if got == v {
			return true
		}
	}
	return false
}

// listableGVRs walks the discovery API and returns a map of group-qualified
// resource name (as produced by qualifiedResourceName) to its GVR, for every
// listable, non-subresource type. The first version seen for a name wins.
// Partial discovery results are used even when discovery reports an error.
func listableGVRs(disco discovery.DiscoveryInterface) (map[string]schema.GroupVersionResource, error) {
	_, lists, err := disco.ServerGroupsAndResources()
	if len(lists) == 0 && err != nil {
		return nil, err
	}

	out := make(map[string]schema.GroupVersionResource)
	for _, list := range lists {
		gv, parseErr := schema.ParseGroupVersion(list.GroupVersion)
		if parseErr != nil {
			continue
		}
		for _, ar := range list.APIResources {
			if strings.Contains(ar.Name, "/") || !hasVerb(ar.Verbs, "list") {
				continue
			}
			name := qualifiedResourceName(ar.Name, gv.Group)
			if _, ok := out[name]; !ok {
				out[name] = gv.WithResource(ar.Name)
			}
		}
	}
	return out, nil
}

// probePopulated checks each GVR for at least one object cluster-wide (a Limit:1
// list, which for namespaced types spans all namespaces). It returns a map keyed
// only by names that were checked successfully; the value is true when the type
// has >=1 object. Names whose check errored (RBAC, timeout) are absent, so
// callers can distinguish "known empty" (present, false) from "unknown" (absent).
func probePopulated(ctx context.Context, dyn dynamic.Interface, gvrByName map[string]schema.GroupVersionResource) map[string]bool {
	var (
		mu        sync.Mutex
		wg        sync.WaitGroup
		populated = make(map[string]bool)
		sem       = make(chan struct{}, probeConcurrency)
	)

	for name, gvr := range gvrByName {
		wg.Add(1)
		sem <- struct{}{}
		go func(name string, gvr schema.GroupVersionResource) {
			defer wg.Done()
			defer func() { <-sem }()

			list, err := dyn.Resource(gvr).List(ctx, metav1.ListOptions{Limit: 1})
			if err != nil {
				return
			}
			mu.Lock()
			populated[name] = len(list.Items) > 0
			mu.Unlock()
		}(name, gvr)
	}

	wg.Wait()
	return populated
}

// probePopulatedStream is the streaming counterpart to probePopulated: it dispatches
// existence checks in `order` (so higher-priority names are checked first) with the
// same bounded concurrency, and calls onPopulated(name) — possibly from several
// goroutines — as soon as a type is found to have >=1 object. It returns when every
// probe has finished or ctx is done. Failures are silently skipped.
func probePopulatedStream(
	ctx context.Context,
	dyn dynamic.Interface,
	gvrByName map[string]schema.GroupVersionResource,
	order []string,
	onPopulated func(name string),
) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, probeConcurrency)

	for _, name := range order {
		if ctx.Err() != nil {
			break
		}
		gvr, ok := gvrByName[name]
		if !ok {
			continue
		}
		sem <- struct{}{}
		wg.Add(1)
		go func(name string, gvr schema.GroupVersionResource) {
			defer wg.Done()
			defer func() { <-sem }()

			list, err := dyn.Resource(gvr).List(ctx, metav1.ListOptions{Limit: 1})
			if err == nil && len(list.Items) > 0 {
				onPopulated(name)
			}
		}(name, gvr)
	}

	wg.Wait()
}
