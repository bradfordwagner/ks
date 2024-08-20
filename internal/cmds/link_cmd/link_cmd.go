package link_cmd

import (
	"errors"
	"fmt"
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/bradfordwagner/ks/internal/link"
	"github.com/bradfordwagner/ks/internal/list"
	"github.com/koki-develop/go-fzf"
)

func Run(a args.PipeArgs) (err error) {
	l := log.Log()

	// list kubeconfigs
	configs, err := list.Kubeconfigs(a.Directory)
	if err != nil {
		l.With("error", err).Error("error listing kubeconfigs")
		return
	}

	// choose a kubeconfig
	one, err := choose.One(configs)
	if errors.Is(fzf.ErrAbort, err) {
		return nil
	} else if err != nil {
		l.With("error", err).Error("error choosing kubeconfig")
		return
	}
	fmt.Print(one)

	// link the chosen kubeconfig
	source := fmt.Sprintf("%s/%s", a.Directory, one)
	target := fmt.Sprintf("%s/config", a.Directory)
	err = link.ForceLink(source, target)
	if err != nil {
		l.With("error", err).Error("error linking kubeconfig")
	}

	return
}
