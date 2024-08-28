package tmux

import (
	"fmt"
	"os/exec"
)

func NewWindow(dir, kubeconfig string) (err error) {
	resolvedKubeconfig := fmt.Sprintf("%s/%s", dir, kubeconfig)
	kubeconfigEnv := fmt.Sprintf("KUBECONFIG=%s", resolvedKubeconfig)
	newWindow := exec.Command("tmux", "new-window", "-n", kubeconfig, "-e", kubeconfigEnv)
	err = newWindow.Run()
	if err != nil {
		return
	}

	return LoadBufferKubeconfig(resolvedKubeconfig)
}
