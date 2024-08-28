package tmux

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
)

// LoadBufferKubeconfig loads the kubeconfig into a tmux buffer
func LoadBufferKubeconfig(kubeconfig string) (err error) {
	//echo ${buffer} | tmux loadb -b ${category}-${index} -
	index := strconv.Itoa(rand.Int())
	bufferName := fmt.Sprintf("kube-%s", index)
	cmd := exec.Command("tmux", "loadb", "-b", bufferName, "-")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("export KUBECONFIG=%s\n", kubeconfig))
	return cmd.Run()
}
