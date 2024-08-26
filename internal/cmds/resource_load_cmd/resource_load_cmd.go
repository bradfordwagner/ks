package resource_load_cmd

import (
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/kube"
	"github.com/bradfordwagner/ks/internal/resources"
	"sort"
	"strings"
)

func Run(a args.Standard) (err error) {
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
	for _, resource := range kuberesources {
		for _, apiResource := range resource.APIResources {
			// exclude sub resource types eg:
			// prioritylevelconfigurations/status
			// cronjobs/status
			if !strings.Contains(apiResource.Name, "/") {
				r.Names = append(r.Names, apiResource.Name)
			}
		}
	}
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
