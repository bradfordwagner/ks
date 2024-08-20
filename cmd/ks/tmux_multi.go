package main

import (
	"fmt"
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/cmds/tmux_multi_cmd"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	fs := tmuxMultiCmd.Flags()
	home, _ := os.UserHomeDir()
	flag_helper.CreateFlag(fs, &tmuxMultiArgs.Directory, "dir", "d", fmt.Sprintf("%s/.kube", home), "env.KS_DIR,default=home/.kube")
}

var tmuxMultiArgs args.TmuxMultiArgs

var tmuxMultiCmd = &cobra.Command{
	Use:   "tmux_multi",
	Short: "creates a new pane in tmux for each kubeconfig",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&tmuxMultiArgs)
		return tmux_multi_cmd.Run(tmuxMultiArgs)
	},
}
