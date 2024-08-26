package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(resourceCmd.Flags())
}

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "resource opens k9s with selected resource view",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.Resource(standardArgs)
	},
}
