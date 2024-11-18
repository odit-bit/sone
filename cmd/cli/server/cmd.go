package server

import (
	"github.com/spf13/cobra"
)

func init() {
	RunCMD.Flags().Int("httpPort", 9696, "listen http port")
	RunCMD.Flags().Int("rtmpPort", 1935, "listen http port")
	RunCMD.Flags().Bool("debug", false, "print log to stdout/stderr")

	ServerCMD.AddCommand(&RunCMD)
}

var ServerCMD = cobra.Command{
	Use: "server",
}

var RunCMD = cobra.Command{
	Use:        "run dir",
	Aliases:    []string{},
	SuggestFor: []string{},
	Short:      "",
	GroupID:    "",
	Long:       "",
	Example:    "",
	ValidArgs:  []string{},

	Args:    cobra.MinimumNArgs(1),
	Version: "0.1",
	Run: func(cmd *cobra.Command, args []string) {
		httpPort, _ := cmd.Flags().GetInt("httpPort")
		rtmpPort, _ := cmd.Flags().GetInt("rtmpPort")
		debug, _ := cmd.Flags().GetBool("debug")

		startServer(args[0], httpPort, rtmpPort, debug)

	},
}
