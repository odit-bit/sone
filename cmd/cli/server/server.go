package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/odit-bit/sone/ingress"
	"github.com/odit-bit/sone/pkg/kvstore"
	"github.com/odit-bit/sone/pkg/rtmp"
	"github.com/odit-bit/sone/pkg/tcp"
	"github.com/odit-bit/sone/streaming"
	"github.com/odit-bit/sone/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func startServer(dir string, port, rtmpPort int, debug bool) {
	// var dir string
	// var port int
	// var rtmpPort int
	// var debug bool
	// flag.StringVar(&dir, "dir", "", "dir to store / cache stream media")
	// flag.IntVar(&port, "http", 9696, "port to listen http")
	// flag.IntVar(&rtmpPort, "rtmp", 1935, "port to listen rtmp ingress")
	// flag.BoolVar(&debug, "debug", false, "print log to terminal")
	// flag.Parse()

	if dir == "" {
		panic("path to dir media empty, memory filesystem still BUG")
	}

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)

	// setup logfile
	f, err := setupLogFile(dir)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.Println("log file:", f.Name())

	//setup logger
	var logWriter io.Writer
	if debug {
		logWriter = io.MultiWriter(f, os.Stderr)
	} else {
		logWriter = f
	}
	logger := logrus.StandardLogger()
	logger.SetOutput(logWriter)

	//setup filesystem
	afs := initFileSystem(dir)
	////////

	//setup errgroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//RTMP listener
	rtmpAddr := fmt.Sprintf(":%v", rtmpPort)
	rtmpL, err := tcp.Listen(ctx, rtmpAddr)
	if err != nil {
		panic(err)
	}

	//HLS listener
	httpAddr := fmt.Sprintf(":%d", port)
	httpL, err := tcp.Listen(ctx, httpAddr)
	if err != nil {
		panic(err)
	}

	//setup kvStore
	kv := kvstore.Open()

	// // HLS Stream service
	// mux := chi.NewRouter()
	// media.InitModule(logger, afs, mux, kv)
	// srv := http.Server{Handler: mux}
	// go srv.Serve(httpL)
	// logger.Printf("listen hls-service on : %v", httpL.Addr().String())

	mux := chi.NewMux()
	// Streaming service
	streaming.StartUp(logger, afs, kv, mux)
	// WebServer (ui)
	web.StartUp(httpL, mux)
	// http-server
	srv := http.Server{Handler: mux}
	go srv.Serve(httpL)
	logger.Printf("listen streaming on : %v", httpL.Addr().String())

	// RTMP ingress service
	h := rtmp.NewHandler()
	ingress.InitModule(ctx, dir, afs, kv, h)
	ingressSrv := rtmp.Server{Handler: h}
	go ingressSrv.Serve(rtmpL)
	logger.Printf("listen ingress on : %v", rtmpL.Addr().String())

	// gracefull shutdown
	select {
	case s := <-sigC:
		logger.Println("got signal", s)
		err = ctx.Err()
	case <-ctx.Done():
		logger.Println("ctx done")
	}
	err = errors.Join(err, srv.Close(), ingressSrv.Close(), rtmpL.Close(), httpL.Close())
	logger.Errorf("shutdown server %v ", err)
}

func setupLogFile(parent string) (*os.File, error) {
	logFilePath, err := filepath.Abs(filepath.Join(parent, "log")) //os.MkdirAll(,0775)
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(logFilePath, 0666); err != nil {
		return nil, err
	}
	filename := filepath.Join(logFilePath, time.Now().Local().Format(time.DateOnly))
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
}

func initFileSystem(dir string) afero.Fs {
	afs := afero.NewOsFs()
	if dir == "" {
		mFs := afero.NewMemMapFs()
		afs = mFs
	} else {
		afs = afero.NewBasePathFs(afs, dir)
	}

	return afs
}
