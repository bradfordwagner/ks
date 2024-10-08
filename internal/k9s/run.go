package k9s

import (
	"os"
	"os/exec"
)

// Run executes k9s with args with interactive terminal
func Run(args ...string) {
	a := append([]string{"--headless"}, args...)
	cmd := exec.Command("k9s", a...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
