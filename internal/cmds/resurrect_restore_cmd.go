package cmds

import (
	"fmt"

	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/resurrect"
)

func ResurrectRestore(a args.Standard) error {
	state, err := resurrect.Load(a.DataDir)
	if err != nil || len(state.Panes) == 0 {
		return err
	}

	panes, err := resurrect.ListPanes()
	if err != nil {
		return err
	}

	// Index live panes by stable positional key.
	type key struct {
		session   string
		windowIdx int
		paneIdx   int
	}
	liveByPos := make(map[key]resurrect.TmuxPane, len(panes))
	for _, p := range panes {
		liveByPos[key{p.Session, p.WindowIdx, p.PaneIdx}] = p
	}

	for _, saved := range state.Panes {
		live, ok := liveByPos[key{saved.Session, saved.WindowIdx, saved.PaneIdx}]
		if !ok {
			continue
		}

		if saved.Kubeconfig != "" {
			_ = resurrect.SendKeys(live.PaneID, fmt.Sprintf("export KUBECONFIG=%s", saved.Kubeconfig))
		}

		if saved.Resource != "" {
			verb := saved.Verb
			if verb == "" {
				verb = "resource"
			}
			_ = resurrect.SendKeys(live.PaneID, fmt.Sprintf("KS_RESOURCE=%s ks %s", saved.Resource, verb))
		}
	}
	return nil
}
