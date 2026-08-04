package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

var resourceLiveAll bool

func init() {
	standardFlagsInit(resourceLiveCmd.Flags())
	resourceLiveCmd.Flags().BoolVarP(&resourceLiveAll, "all", "A", false, "watch selections across all namespaces")
}

var resourceLiveCmd = &cobra.Command{
	Use:   "resource_live",
	Short: "live-scan the cluster and pick from resource types that currently have objects",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.ResourceLive(standardArgs, resourceLiveAll)
	},
}
