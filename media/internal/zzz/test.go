package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/odit-bit/sone/media/internal/application"
	mediaGRPC "github.com/odit-bit/sone/media/internal/grpc"
	"github.com/odit-bit/sone/media/internal/segment"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGINT)

	g := errgroup.Group{}
	memfs := afero.NewMemMapFs()
	repo := segment.NewLocalFS(memfs)
	md := application.New(repo)

	srv := grpc.NewServer()
	if err := mediaGRPC.RegisterServer(md, srv); err != nil {
		log.Println(err)
		return
	}

	l, err := net.Listen("tcp", "localhost:6969")
	if err != nil {
		panic(err)
	}

	g.Go(func() error {
		<-sigC
		srv.Stop()
		return nil
	})

	g.Go(func() error {
		return srv.Serve(l)
	})

	g.Wait()
	fmt.Println("close media server")
}
