package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(clearCacheCmd.Flags())
}

var clearCacheCmd = &cobra.Command{
	Use:   "clear_cache",
	Short: "clears cache tmux resource panes",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.ClearCache(standardArgs, true)
	},
}
