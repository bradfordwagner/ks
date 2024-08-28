package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(kubeCopyCommand.Flags())
}

var kubeCopyCommand = &cobra.Command{
	Use:   "kube_cp",
	Short: "copies kubecontext to tmux buffer",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.KubeCopy(standardArgs)
	},
}
