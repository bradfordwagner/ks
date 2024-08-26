package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(tmuxMultiCmd.Flags())
}

var tmuxMultiCmd = &cobra.Command{
	Use:   "tmux_multi",
	Short: "creates a new pane in tmux for each kubeconfig",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.TmuxMulti(standardArgs)
	},
}
