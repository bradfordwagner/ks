package main

import (
	"fmt"
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/cmds/pipe_cmd"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	fs := pipeCmd.Flags()
	home, _ := os.UserHomeDir()
	flag_helper.CreateFlag(fs, &pipeArgs.Directory, "dir", "d", fmt.Sprintf("%s/.kube", home), "env.KS_DIR,default=home/.kube")
}

var pipeArgs args.PipeArgs

var pipeCmd = &cobra.Command{
	Use:   "pipe",
	Short: "pipes kubecontext selection to stdout",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flag_helper.Load(&pipeArgs)
		return pipe_cmd.Run(pipeArgs)
	},
}
