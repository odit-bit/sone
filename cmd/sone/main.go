package main

import (
	"github.com/odit-bit/sone/cmd/cli/client"
	"github.com/odit-bit/sone/cmd/cli/server"
	"github.com/spf13/cobra"
)

func main() {

	cmd := cobra.Command{}
	cmd.AddCommand(&server.ServerCMD)
	cmd.AddCommand(&client.ClientCMD)
	cmd.Execute()
}
