package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/glebarez/go-sqlite"
	streamGRPC "github.com/odit-bit/sone/streaming/grpc"
	"github.com/odit-bit/sone/streaming/streamingpb"
	"golang.org/x/sync/errgroup"

	"github.com/odit-bit/sone/streaming/internal/api"
	"github.com/odit-bit/sone/streaming/internal/app"
	"github.com/odit-bit/sone/streaming/internal/database"
	"google.golang.org/grpc"
)

func main() {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, os.Interrupt)
	g := errgroup.Group{}
	g.SetLimit(3)

	l, err := net.Listen("tcp", "localhost:6969")
	if err != nil {
		panic(err)
	}

	dsn := ":memory:"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatal(err)
		return
	}

	//construct application neccessary component.
	repo, err := database.NewSqlite(db)
	if err != nil {
		panic(err)
	}
	rpc := api.NewMockGRPC()
	app := app.New(&repo, &rpc)
	_ = app

	///////
	srv := grpc.NewServer()
	ss := streamGRPC.NewServer(app)
	streamingpb.RegisterLiveStreamServer(srv, ss)

	g.Go(func() error {
		<-sigC
		fmt.Println("got signal")
		srv.Stop()
		l.Close()
		return nil
	})

	g.Go(func() error {
		return srv.Serve(l)
	})

	g.Wait()
}
