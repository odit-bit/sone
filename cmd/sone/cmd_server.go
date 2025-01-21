package main

import (
	"github.com/odit-bit/sone/cmd/app"
	"github.com/spf13/cobra"
)

func init() {
	RunCMD.Flags().Int("httpPort", 9797, "listen http port")
	RunCMD.Flags().Int("rpcPort", 9696, "listen rpc port")
	RunCMD.Flags().Int("rtmpPort", 1935, "listen rtmp port")
	RunCMD.Flags().Bool("debug", false, "print log to stdout/stderr")
	RunCMD.Flags().String("fs", "", "path to dir")

	RunCMD.Flags().String("minio-url", "", "minio address")
	RunCMD.Flags().String("minio-key", "", "minio access key")
	RunCMD.Flags().String("minio-secret", "", "minio secret key")

	RunCMD.Flags().String("sql", "", "sql dsn string")

	ServerCMD.AddCommand(&RunCMD)
}

var ServerCMD = cobra.Command{
	Use: "server",
}

var RunCMD = cobra.Command{
	Use:        "run",
	Aliases:    []string{},
	SuggestFor: []string{},
	Short:      "",
	GroupID:    "",
	Long:       "",
	Example:    "",
	ValidArgs:  []string{},

	Args:    cobra.MaximumNArgs(0),
	Version: "0.1",
	Run: func(cmd *cobra.Command, args []string) {
		httpPort, _ := cmd.Flags().GetInt("httpPort")
		rpcPort, _ := cmd.Flags().GetInt("rpcPort")
		rtmpPort, _ := cmd.Flags().GetInt("rtmpPort")
		debug, _ := cmd.Flags().GetBool("debug")

		fsPath, _ := cmd.Flags().GetString("fs")

		minioUrl, _ := cmd.Flags().GetString("minio-url")
		minioKey, _ := cmd.Flags().GetString("minio-key")
		minioSecret, _ := cmd.Flags().GetString("minio-secret")

		sqlString, _ := cmd.Flags().GetString("sql")

		app.Start(app.Config{
			Http: struct{ Port int }{
				httpPort,
			},
			Rpc: struct{ Port int }{
				rpcPort,
			},
			Rtmp: struct{ Port int }{
				rtmpPort,
			},

			Logging: struct{ Debug bool }{
				debug,
			},

			Minio: struct {
				Address   string
				AccessKey string
				SecretKey string
			}{
				Address:   minioUrl,
				AccessKey: minioKey,
				SecretKey: minioSecret,
			},

			Filesystem: struct{ Path string }{
				fsPath,
			},

			SQL: struct{ DSN string }{
				DSN: sqlString,
			},
		})

	},
}
