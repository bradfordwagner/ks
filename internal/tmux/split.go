package tmux

import (
	"fmt"
	"github.com/bradfordwagner/go-util/log"
	"os/exec"
	"strings"
)

// splitTiled sets the tiled layout, splits a detached pane with KUBECONFIG set,
// and returns the new pane index.
func splitTiled(kubeconfig string) (paneIndex string, err error) {
	l := log.Log()

	err = TiledLayout()
	if err != nil {
		return "", err
	}

	kubeconfigEnv := fmt.Sprintf("KUBECONFIG=%s", kubeconfig)
	splitw := exec.Command("tmux", "splitw", "-d", "-P", "-F", "#{pane_index}", "-e", kubeconfigEnv)
	output, err := splitw.Output()
	if err != nil {
		l.With("error", err).Error("couldn't split window")
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Split opens a new tiled pane with KUBECONFIG set and echoes the kubeconfig path.
func Split(kubeconfig string) (err error) {
	paneIndex, err := splitTiled(kubeconfig)
	if err != nil {
		return err
	}
	return SendKeys(paneIndex, "echo kube_config=${KUBECONFIG}")
}

// SplitAndSend opens a new tiled pane with KUBECONFIG set and sends keys to it.
func SplitAndSend(kubeconfig, keys string) (err error) {
	paneIndex, err := splitTiled(kubeconfig)
	if err != nil {
		return err
	}
	return SendKeys(paneIndex, keys)
}
