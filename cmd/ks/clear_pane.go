package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(clearPaneCmd.Flags())
}

var clearPaneCmd = &cobra.Command{
	Use:   "clear_pane",
	Short: "clears tmux pane cache resource",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.ClearCache(standardArgs, false)
	},
}
