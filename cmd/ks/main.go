package main

import (
	"fmt"
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use: "ks",
}

func init() {
	rootCmd.AddCommand(
		linkCmd,
		pipeCmd,
		resourceCmd,
		resourceLoadCommand,
		tmuxMultiCmd,
		tmuxWindowCmd,
	)
}

func standardFlagsInit(fs *pflag.FlagSet) {
	home, _ := os.UserHomeDir()
	flag_helper.CreateFlag(fs, &standardArgs.Directory, "dir", "d", fmt.Sprintf("%s/.kube", home), "env.KS_DIR,default=home/.kube")
	flag_helper.CreateFlag(fs, &standardArgs.Kubeconfig, "kubeconfig", "k", fmt.Sprintf("%s/.kube/config", home), "env.KUBECONFIG,default=home/.kube/config")
	flag_helper.CreateFlag(fs, &standardArgs.Timeout, "timeout", "t", 10*time.Second, "default=10s")
}

var standardArgs args.Standard

func main() {
	l := log.Log()
	// cobra
	if err := rootCmd.Execute(); err != nil {
		l.With("error", err).Fatal("could not execute command")
	}
}
