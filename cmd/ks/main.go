package main

import (
	"errors"
	"fmt"
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/koki-develop/go-fzf"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"time"
)

var commands []*cobra.Command

func init() {
	commands = []*cobra.Command{
		resourceCmd,
		resourceAllCmd,
		setNamespaceCmd,
		tmuxMultiCmd,
		tmuxWindowCmd,
		kubeCopyCommand,
		linkCmd,
		pipeCmd,
		resourceLoadCommand,
	}
	rootCmd.AddCommand(commands...)
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

var rootCmd = &cobra.Command{
	Use: "ks",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// map command names to commands
		nameToCommand := make(map[string]*cobra.Command)
		var names []string
		for _, command := range commands {
			nameToCommand[command.Name()] = command
			names = append(names, command.Name())
		}

		// choose a command
		selectedCommand, err := choose.One(names)
		if errors.Is(fzf.ErrAbort, err) {
			return nil
		} else if err != nil {
			return err
		}

		return nameToCommand[selectedCommand].RunE(cmd, args)
	},
}
