package tmux

import (
	"fmt"
	"os/exec"
)

func NewWindow(dir, kubeconfig string) (err error) {
	kubeconfigEnv := fmt.Sprintf("KUBECONFIG=%s/%s", dir, kubeconfig)
	newWindow := exec.Command("tmux", "new-window", "-n", kubeconfig, "-e", kubeconfigEnv)
	err = newWindow.Run()
	if err != nil {
		return
	}

	return SendKeys("0", "kcc")
}
