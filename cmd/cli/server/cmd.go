package server

import (
	"github.com/spf13/cobra"
)

func init() {
	RunCMD.Flags().Int("httpPort", 9797, "listen http port")
	RunCMD.Flags().Int("rpcPort", 9696, "listen rpc port")
	RunCMD.Flags().Int("rtmpPort", 1935, "listen rtmp port")
	RunCMD.Flags().Bool("debug", false, "print log to stdout/stderr")
	RunCMD.Flags().Bool("verbose", false, "set level to debug level")
	RunCMD.Flags().Bool("test", false, "set deployment mode")

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
		rpcPort, _ := cmd.Flags().GetInt("rpcPort")
		rtmpPort, _ := cmd.Flags().GetInt("rtmpPort")
		debug, _ := cmd.Flags().GetBool("debug")
		verbose, _ := cmd.Flags().GetBool("verbose")
		isTest, _ := cmd.Flags().GetBool("test")

		startApp(args[0], rpcPort, rtmpPort, httpPort, debug, verbose, isTest)

	},
}
