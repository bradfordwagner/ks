package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(resourceLiveAllCmd.Flags())
}

var resourceLiveAllCmd = &cobra.Command{
	Use:   "resource_live_all",
	Short: "live-scan the cluster and pick from resource types that currently have objects, across all namespaces",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.ResourceLive(standardArgs, true)
	},
}
