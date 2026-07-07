package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(resurrectRestoreCmd.Flags())
}

var resurrectRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "tmux-resurrect post-restore hook: re-apply KUBECONFIG and relaunch ks resource per pane",
	RunE: func(cmd *cobra.Command, args []string) error {
		flag_helper.Load(&standardArgs)
		return cmds.ResurrectRestore(standardArgs)
	},
}
