package main

import (
	"github.com/odit-bit/sone/cmd/client"
	"github.com/spf13/cobra"
)

func main() {

	cmd := cobra.Command{}
	cmd.AddCommand(&ServerCMD)
	cmd.AddCommand(&client.ClientCMD)
	cmd.Execute()
}
