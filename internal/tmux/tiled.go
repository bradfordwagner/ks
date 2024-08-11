package tmux

import "os/exec"

func TiledLayout() (err error) {
	setLayout := exec.Command("tmux", "select-layout", "tiled")
	return setLayout.Run()
}
