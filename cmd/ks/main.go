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
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	setEnv bool // global flag to tell if we are using KUBECONFIG
	tmux bool
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
	Run: func(cmd *cobra.Command, args []string) {
		execute(true, args)
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
	Run: func(cmd *cobra.Command, args []string) {
		execute(false, args)
	},
}

func init() {
	rootCmd.AddCommand(localCmd, kubeCmd)
	rootCmd.PersistentFlags().BoolVarP(&setEnv, "setenv", "s", false, "copies export KUBECONTEXT")
	rootCmd.PersistentFlags().BoolVarP(&tmux, "tmux", "t", false, "executes export KUBECONTEXT in a new tmux pane")
}

func extractKubeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	return fmt.Sprintf("%s/.kube", homeDir)
}

func execute(isLocal bool, args []string) {
	maxLength := 1
	if tmux {
		maxLength = len(args)
	}

	for i := 0; i < maxLength; i++ {
		filePath := args[i]
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

		if tmux {
			tmuxSplit(filePath)
		}

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
}

func tmuxSplit(path string) {
	setLayout := exec.Command("tmux", "select-layout", "tiled")
	setLayout.Run()

	splitw := exec.Command("tmux", "splitw", "-d", "-P", "-F", "#{pane_index}")
	output, err := splitw.Output()
	if err != nil {
		logrus.WithError(err).Error("couldn't run hello world")
	}
	paneIndex := strings.TrimSpace(string(output))

	tmuxSendToPane(paneIndex, fmt.Sprintf("export KUBECONFIG=%s", path))

	// this is a bradford special.. should it really be here?
	tmuxSendToPane(paneIndex, "kcompletion")
}

func tmuxSendToPane(paneIndex, command string)  {
	setKubeConfig := exec.Command("tmux", "send", "-t", paneIndex, command, "ENTER")
	_ = setKubeConfig.Start()
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
