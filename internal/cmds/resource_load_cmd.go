package cmds

import (
	"github.com/bradfordwagner/go-util/bwutil"
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/kube"
	"github.com/bradfordwagner/ks/internal/resources"
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

	var r resources.Resources
	resourceNames := bwutil.NewSet[string]()
	for _, resource := range kuberesources {
		for _, apiResource := range resource.APIResources {
			// exclude sub resource types eg:
			// prioritylevelconfigurations/status
			// cronjobs/status
			if !strings.Contains(apiResource.Name, "/") {
				resourceNames.Add(apiResource.Name)
			}
		}
	}
	r.Names = resourceNames.Keyset()
	sort.Strings(r.Names)

	for _, n := range r.Names {
		l.With("resource", n).Info("found resource")
	}

	resourceFile, err := r.Write(a.Directory)
	if err != nil {
		return err
	}
	l.With("file", resourceFile).Info("wrote resources to file")

	return
}
