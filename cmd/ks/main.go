package main

import (
	"github.com/bradfordwagner/go-util/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "ks",
}

func init() {
	rootCmd.AddCommand(
		linkCmd,
		pipeCmd,
		tmuxMultiCmd,
		tmuxWindowCmd,
	)
}

func main() {
	l := log.Log()
	// cobra
	if err := rootCmd.Execute(); err != nil {
		l.With("error", err).Fatal("could not execute command")
	}
}
