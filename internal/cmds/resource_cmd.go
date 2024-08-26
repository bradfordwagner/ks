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
func Resource(a args.Standard) (err error) {
	// Load resources
	r, err := resources.LoadResources(a.Directory)
	if err != nil {
		return
	}

	// choose a resource type
	resourceType, err := choose.One(r.Names)
	if errors.Is(err, fzf.ErrAbort) {
		return nil
	} else if err != nil {
		return
	}

	// run k9s with the selected resource type
	k9s.Run("-c", resourceType)

	return
}
