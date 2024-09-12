package cmds

import (
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/resources"
)

func ClearCache(a args.Standard, all bool) (err error) {
	l := log.Log()

	// Load resources
	r, err := resources.LoadResources(a.Directory)
	if err != nil {
		l.With("error", err).Error("failed to load resources")
		return err
	}

	if all {
		l.Info("clearing all cache")
		r.ResetCache()
	} else {
		l.Info("clearing pane cache")
		r.ResetPane()
	}
	_, err = r.Write(a.Directory)
	return
}
