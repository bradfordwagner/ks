package tmux

import "os/exec"

func SendKeys(paneIndex, command string) (err error) {
	setKubeConfig := exec.Command("tmux", "send", "-t", paneIndex, command, "ENTER")
	return setKubeConfig.Run()
}
