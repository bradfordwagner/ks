package cmds

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

func KubeCopy(a args.Standard) (err error) {
	l := log.Log()

	// list kubeconfigs
	configs, err := list.Kubeconfigs(a.Directory)
	if err != nil {
		l.With("error", err).Error("error listing kubeconfigs")
		return
	}

	// choose a kubeconfig
	kubeconfig, err := choose.One(configs)
	if errors.Is(fzf.ErrAbort, err) {
		return nil
	} else if err != nil {
		l.With("error", err).Error("error choosing kubeconfig")
		return
	}
	kubeconfig = fmt.Sprintf("%s/%s", a.Directory, kubeconfig)
	l.With("kubeconfig", kubeconfig).Info("kubeconfig copied to tmux buffer")

	// copy to tmux buffer
	return tmux.LoadBufferKubeconfig(kubeconfig)
}
