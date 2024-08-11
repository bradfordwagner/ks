package main

import (
	"fmt"
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds/tmux_window_cmd"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	fs := tmuxWindowCmd.Flags()
	home, _ := os.UserHomeDir()
	flag_helper.CreateFlag(fs, &tmuxMultiArgs.Directory, "dir", "d", fmt.Sprintf("%s/.kube", home), "env.KS_DIR,default=home/.kube")
}

var tmuxWindowCmd = &cobra.Command{
	Use:   "tmux_window",
	Short: "creates a new window for selected kubeconfig and copies it to tmux buffer",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&tmuxMultiArgs)
		return tmux_window_cmd.Run(tmuxMultiArgs)
	},
}
