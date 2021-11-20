//usr/bin/env go run "$0" "$@"; exit "$?"

/*
requires:
brew install pbcopy
*/

package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	setEnv bool // global flag to tell if we are using KUBECONFIG
)

var rootCmd = &cobra.Command{
	Use: "ks",
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "sets a local file to be the kubecontext",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{""}, cobra.ShellCompDirectiveFilterFileExt
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		execute(true, args[0])
	},
}

var kubeCmd = &cobra.Command{
	Use:   "kube",
	Short: "sets a file from ~/.kube to be kubecontext",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		kubeDirPattern := extractKubeDir()
		files, err := ioutil.ReadDir(kubeDirPattern)
		if err != nil {
			logrus.WithError(err).Fatal()
		}
		filesInKubeDir := []string{}
		for _, file := range files {
			filesInKubeDir = append(filesInKubeDir, file.Name())
		}
		return filesInKubeDir, cobra.ShellCompDirectiveDefault
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		execute(false, args[0])
	},
}

func init() {
	rootCmd.AddCommand(localCmd, kubeCmd)
	rootCmd.PersistentFlags().BoolVarP(&setEnv, "setenv", "s", false, "copies export KUBECONTEXT")
}

func extractKubeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	return fmt.Sprintf("%s/.kube", homeDir)
}

func execute(isLocal bool, filePath string) {
	// local flag was passed in
	if isLocal {
		// find absolute path to provided relative file path
		if p, err := filepath.Abs(filePath); err != nil {
			logrus.WithError(err).Fatalf("could not find file=%s", filePath)
		} else {
			filePath = p
		}
	} else {
		// use ~/.kube/... path for kube verb
		filePath = fmt.Sprintf("%s/%s", extractKubeDir(), filePath)
	}
	logrus.Infof("context file=%s", filePath)

	if setEnv {
		if err := execKubeContextCommand(filePath); err != nil {
			logrus.WithError(err).Fatal("could not copy export command")
		} else {
			logrus.Info("copied kubeconfig command to clipboard, paste when ready")
			os.Exit(0)
		}
	}

	// link .kube file with filePath
	kubeConfigPath := fmt.Sprintf("%s/config", extractKubeDir())
	_ = os.RemoveAll(kubeConfigPath) // we don't care if this file exists or not
	if err := os.Symlink(filePath, kubeConfigPath); err != nil {
		logrus.WithError(err).Fatal("could not link configurations")
	}
}

func execKubeContextCommand(filePath string) (err error) {
	copyText := fmt.Sprintf("export KUBECONFIG=%s", filePath)
	return clipboard.WriteAll(copyText)
}

func main() {
	// setup logrus
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// cobra
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
