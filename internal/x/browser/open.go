package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

var (
	url = flag.String("url", "localhost:9696", "open url in browser")
)

func main() {
	flag.Parse()
	err := openURL(*url)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
}

func openURL(url string) error {
	var err error
	fmt.Println("opening url", url)
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("cannot open URL %s on this platform", url)
	}
	return err
}
