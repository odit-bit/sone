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
	"github.com/odit-bit/sone/users/internal/application"
	userGRPC "github.com/odit-bit/sone/users/internal/grpc"
	"github.com/odit-bit/sone/users/internal/sqlite"
	"github.com/odit-bit/sone/users/userpb"
	"golang.org/x/sync/errgroup"
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

	users, err := sqlite.New("users", db)
	if err != nil {
		panic(err)
	}

	// // streams := kv.New(mono.KV())
	// glusers, err := sqlite.NewUserGlRepo(context.Background(), db, "glusers")
	// if err != nil {
	// 	mono.Logger().Panic(err)
	// }

	///////
	srv := grpc.NewServer()
	app := application.New(nil, users, nil)
	u := userGRPC.New(app)
	userpb.RegisterUserServiceServer(srv, u)

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
