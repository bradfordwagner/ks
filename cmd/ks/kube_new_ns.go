package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

func init() {
	standardFlagsInit(kubeNewNamespaceCommand.Flags())
}

var kubeNewNamespaceCommand = &cobra.Command{
	Use:   "kube_new_ns",
	Short: "creates a new kubeconfig with namespace appended, and copies it to tmux buffer",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&standardArgs)
		return cmds.KubeNewNamespace(standardArgs)
	},
}
