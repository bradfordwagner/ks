package tmux_window_cmd

import (
	"errors"
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/bradfordwagner/ks/internal/list"
	"github.com/bradfordwagner/ks/internal/tmux"
	"github.com/koki-develop/go-fzf"
)

func Run(a args.Standard) (err error) {
	l := log.Log()

	configs, err := list.Kubeconfigs(a.Directory)
	if err != nil {
		l.With("error", err).Error("error listing kubeconfigs")
		return
	}

	one, err := choose.One(configs)
	if errors.Is(err, fzf.ErrAbort) {
		return
	} else if err != nil {
		l.With("error", err).Error("error choosing kubeconfig")
		return
	}

	// create new tmux window
	return tmux.NewWindow(a.Directory, one)
}
