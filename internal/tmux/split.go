package tmux

import (
	"fmt"
	"github.com/bradfordwagner/go-util/log"
	"os/exec"
	"strings"
)

func Split(kubeconfig string) (err error) {
	l := log.Log()

	err = TiledLayout()
	if err != nil {
		return err
	}

	kubeconfigEnv := fmt.Sprintf("KUBECONFIG=%s", kubeconfig)
	splitw := exec.Command("tmux", "splitw", "-d", "-P", "-F", "#{pane_index}", "-e", kubeconfigEnv)
	output, err := splitw.Output()
	if err != nil {
		l.With("error", err).Error("couldn't split window")
		return
	}
	paneIndex := strings.TrimSpace(string(output))
	return SendKeys(paneIndex, "echo kube_config=${KUBECONFIG}")
}
