package cmds

import (
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/resurrect"
	"github.com/bradfordwagner/ks/internal/resources"
)

type paneKey struct {
	session   string
	windowIdx int
	paneIdx   int
}

func ResurrectSave(a args.Standard) error {
	panes, err := resurrect.ListPanes()
	if err != nil || len(panes) == 0 {
		return err
	}

	cache, _ := resources.LoadResources(a.DataDir)

	// Load existing sidecar as the source of truth for verb.
	existing, _ := resurrect.Load(a.DataDir)
	verbByPos := make(map[paneKey]string, len(existing.Panes))
	for _, p := range existing.Panes {
		verbByPos[paneKey{p.Session, p.WindowIdx, p.PaneIdx}] = p.Verb
	}

	var state resurrect.ResurrectState
	state.Version = 1

	for _, pane := range panes {
		resource := cache.GetByPane(pane.PaneID)
		env, _ := resurrect.ReadEnv(pane.PanePID)
		kubeconfig := env["KUBECONFIG"]

		if resource == "" && kubeconfig == "" {
			continue
		}

		verb := verbByPos[paneKey{pane.Session, pane.WindowIdx, pane.PaneIdx}]
		if verb == "" {
			verb = verbForPane(pane)
		}

		state.Panes = append(state.Panes, resurrect.ResurrectPane{
			Session:    pane.Session,
			WindowIdx:  pane.WindowIdx,
			PaneIdx:    pane.PaneIdx,
			Kubeconfig: kubeconfig,
			Resource:   resource,
			Verb:       verb,
		})
	}

	if len(state.Panes) == 0 {
		return nil
	}
	return state.Write(a.DataDir)
}

// verbForPane is a fallback when no verb is stored in the sidecar.
// It walks /proc for a k9s descendant and checks its cmdline for -A.
func verbForPane(pane resurrect.TmuxPane) string {
	k9sPID, _ := resurrect.FindDescendantByName(pane.PanePID, "k9s")
	if k9sPID == 0 && pane.CurrentCommand == "k9s" {
		k9sPID = pane.PanePID
	}
	if k9sPID != 0 {
		cmdline, _ := resurrect.ReadCmdline(k9sPID)
		for _, arg := range cmdline {
			if arg == "-A" {
				return "resource_all"
			}
		}
	}
	return "resource"
}
