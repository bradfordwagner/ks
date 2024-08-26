package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(resourceLoadCommand.Flags())
}

var resourceLoadCommand = &cobra.Command{
	Use:   "resources_load",
	Short: "creates a new pane in tmux for each kubeconfig",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.ResourceLoad(standardArgs)
	},
}
