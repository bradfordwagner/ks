package main

import (
	"github.com/bradfordwagner/go-util/flag_helper"
	"github.com/bradfordwagner/ks/internal/cmds"
	"github.com/spf13/cobra"
)

var resourceLeaderboardAll bool

func init() {
	standardFlagsInit(resourceLeaderboardCmd.Flags())
	resourceLeaderboardCmd.Flags().BoolVarP(&resourceLeaderboardAll, "all", "a", false, "show all resources including those with zero votes")
}

var resourceLeaderboardCmd = &cobra.Command{
	Use:   "resource_leaderboard",
	Short: "show resource selection leaderboard sorted by usage",
	RunE: func(cmd *cobra.Command, args []string) error {
		flag_helper.Load(&standardArgs)
		return cmds.ResourceLeaderboard(standardArgs, resourceLeaderboardAll)
	},
}
