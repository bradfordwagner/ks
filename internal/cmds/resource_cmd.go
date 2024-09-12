package cmds

import (
	"errors"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/bradfordwagner/ks/internal/k9s"
	"github.com/bradfordwagner/ks/internal/resources"
	"github.com/koki-develop/go-fzf"
)

// Resource opens k9s with selected resource view
func Resource(a args.Standard, all bool) (err error) {
	// Load resources
	r, err := resources.LoadResources(a.Directory)
	if err != nil {
		return
	}

	// Get the resource type
	resourceType := r.Get()
	if resourceType == "" {
		//choose a resource type
		resourceType, err = choose.One(r.Names)
		if errors.Is(err, fzf.ErrAbort) {
			return nil
		} else if err != nil {
			return
		}

		// save the selected resource type
		r.Upsert(resourceType)
		go r.Write(a.Directory)
	}

	// if all is true, run k9s with all resources
	args := []string{"-c", resourceType}
	if all {
		args = append(args, "-A")
	}

	// run k9s with the selected resource type
	k9s.Run(args...)

	return
}
