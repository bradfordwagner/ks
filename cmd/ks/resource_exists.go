package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

var resourceExistsAll bool

func init() {
	standardFlagsInit(resourceExistsCmd.Flags())
	resourceExistsCmd.Flags().BoolVarP(&resourceExistsAll, "all", "A", false, "watch selections across all namespaces")
}

var resourceExistsCmd = &cobra.Command{
	Use:   "resource_exists",
	Short: "pick from resource types recorded as populated at the last resources_load (cached, no live scan)",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.ResourceExists(standardArgs, resourceExistsAll)
	},
}
