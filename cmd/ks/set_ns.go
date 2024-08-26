package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(setNamespaceCmd.Flags())
}

var setNamespaceCmd = &cobra.Command{
	Use:   "set_ns",
	Short: "set the namespace for the current context",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.SetNamespace(standardArgs)
	},
}
