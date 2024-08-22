package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds/link_cmd"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(linkCmd.Flags())
}

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "links main kubeconfig home/.kube/config to selection",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return link_cmd.Run(standardArgs)
	},
}
