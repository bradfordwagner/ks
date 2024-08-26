package cmds

import (
	"errors"
	"fmt"
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/bradfordwagner/ks/internal/list"
	"github.com/koki-develop/go-fzf"
)

// Pipe is the main entry point for the pipe command
func Pipe(a args.Standard) (err error) {
	l := log.Log()

	configs, err := list.Kubeconfigs(a.Directory)
	if err != nil {
		l.With("error", err).Error("error listing kubeconfigs")
		return
	}

	one, err := choose.One(configs)
	if errors.Is(fzf.ErrAbort, err) {
		return nil
	} else if err != nil {
		l.With("error", err).Error("error choosing kubeconfig")
		return
	}
	res := fmt.Sprintf("%s/%s", a.Directory, one)
	fmt.Print(res)

	return
}
