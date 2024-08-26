package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(pipeCmd.Flags())
}

var pipeCmd = &cobra.Command{
	Use:   "pipe",
	Short: "pipes kubecontext selection to stdout",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.Pipe(standardArgs)
	},
}
