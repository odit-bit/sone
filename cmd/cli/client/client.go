package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	CreateStreamKeyCMD.Flags().String("address", "http://localhost:9696", "sone http server address, default localhost:9696")
	ClientCMD.AddCommand(&CreateStreamKeyCMD)
}

var ClientCMD = cobra.Command{
	Use: "admin",
}

var CreateStreamKeyCMD = cobra.Command{
	Use: "create-stream name",
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("address")
		path := "stream/create"
		u, err := url.Parse(addr)
		if err != nil {
			log.Fatal(err)
		}
		u.Scheme = "http"
		u.Path = path
		req, err := http.NewRequestWithContext(cmd.Context(), "GET", u.String(), nil)
		if err != nil {
			log.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 && res.StatusCode != 202 {
			log.Fatalf("error: %v", res.Status)
		}
		b, _ := io.ReadAll(res.Body)
		fmt.Println("your stream-key:", string(b))
	},
}
