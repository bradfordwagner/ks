package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(tmuxWindowCmd.Flags())
}

var tmuxWindowCmd = &cobra.Command{
	Use:   "tmux_window",
	Short: "creates a new window for selected kubeconfig and copies it to tmux buffer",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.TmuxWindow(standardArgs)
	},
}
