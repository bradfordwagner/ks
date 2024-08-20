package main

import (
	"fmt"
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds/link_cmd"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	fs := linkCmd.Flags()
	home, _ := os.UserHomeDir()
	flag_helper.CreateFlag(fs, &pipeArgs.Directory, "dir", "d", fmt.Sprintf("%s/.kube", home), "env.KS_DIR,default=home/.kube")
}

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "links main kubeconfig home/.kube/config to selection",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&pipeArgs)
		return link_cmd.Run(pipeArgs)
	},
}
