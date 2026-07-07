package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(resurrectSaveCmd.Flags())
}

var resurrectSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "tmux-resurrect pre-save hook: snapshot KUBECONFIG and resource per pane",
	RunE: func(cmd *cobra.Command, args []string) error {
		flag_helper.Load(&standardArgs)
		return cmds.ResurrectSave(standardArgs)
	},
}
