package tmux_multi_cmd

import (
	"errors"
	"fmt"
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

	selected, err := choose.Multi(configs)
	if errors.Is(err, fzf.ErrAbort) {
		return
	} else if err != nil {
		l.With("error", err).Error("error choosing kubeconfig")
		return
	}

	// open a pane per kubeconfig and set the KUBECONFIG env var
	for _, kubeconfig := range selected {
		err := tmux.Split(fmt.Sprintf("%s/%s", a.Directory, kubeconfig))
		if err != nil {
			l.With("error", err).Error("error splitting tmux window")
			return err
		}
	}

	return
}
