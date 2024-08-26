package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(resourceAllCmd.Flags())
}

var resourceAllCmd = &cobra.Command{
	Use:   "resource_all",
	Short: "resource all opens k9s with selected resource view",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.Resource(standardArgs, true)
	},
}
