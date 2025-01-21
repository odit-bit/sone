package app

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/glebarez/go-sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/odit-bit/sone/ingress"
	"github.com/odit-bit/sone/internal/monolith"
	"github.com/odit-bit/sone/media"
	"github.com/odit-bit/sone/pkg/kvstore"
	"github.com/odit-bit/sone/pkg/observer"
	"github.com/odit-bit/sone/pkg/rtmp"
	"github.com/odit-bit/sone/streaming"
	"github.com/odit-bit/sone/users"
	"github.com/odit-bit/sone/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/reflection"
)

var _ monolith.Monolith = (*App)(nil)

type App struct {
	fs          afero.Fs
	kv          *kvstore.Client
	logger      *logrus.Logger
	rpc         *monolith.RPC
	http        *monolith.HTTP
	rtmpHandler *rtmp.HandlerRegister
	db          *sql.DB
	minioConf   *monolith.MinioConfig
	observer    *observer.Observer
}

// FS implements monolith.Monolith.
func (i *App) FS() afero.Fs {
	return i.fs
}

// KV implements monolith.Monolith.
func (i *App) KV() *kvstore.Client {
	return i.kv
}

// Logger implements monolith.Monolith.
func (i *App) Logger() *logrus.Logger {
	return i.logger
}

// RPC implements monolith.Monolith.
func (i *App) RPC() *monolith.RPC {
	return i.rpc
}

// RTMP implements monolith.Monolith
func (i *App) RTMP() *rtmp.HandlerRegister {
	return i.rtmpHandler
}

func (i *App) DB() *sql.DB {
	return i.db
}

func (i *App) Minio() monolith.MinioConfig {
	return *i.minioConf
}

func (i *App) HTTP() *monolith.HTTP {
	return i.http
}

func (i *App) Observer() *observer.Observer {
	return i.observer
}

func (i *App) Run(ctx context.Context, rpcAddr, httpAddr, rtmpAddr string) error {

	eg, ictx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		err := i.runGrpc(ctx, rpcAddr)
		i.logger.Debug(err)
		return err
	})
	eg.Go(func() error {
		err := i.runHTTP(ctx, httpAddr)
		i.logger.Debug(err)
		return err

	})
	eg.Go(func() error {
		err := i.runRTMP(ctx, rtmpAddr)
		i.logger.Debug(err)
		return err
	})

	<-ictx.Done()
	return eg.Wait()

}

func (i *App) runGrpc(ctx context.Context, rpcEndpoint string) error {
	errGroup, ictx := errgroup.WithContext(ctx)
	grpcL, err := net.Listen("tcp", rpcEndpoint)
	if err != nil {
		return err
	}
	errGroup.Go(func() error {
		if err := i.rpc.Serve(grpcL); err != nil {
			return err
		}
		return nil
	})
	i.logger.Info("listen grpc on:", grpcL.Addr().String())
	errGroup.Go(func() error {
		<-ictx.Done()
		i.rpc.Stop()
		return nil
	})
	return errGroup.Wait()
}

func (i *App) runHTTP(ctx context.Context, httpEndpoint string) error {
	errGroup, ictx := errgroup.WithContext(ctx)
	httpL, err := net.Listen("tcp", httpEndpoint)
	if err != nil {
		return err
	}

	srv := http.Server{Handler: i.http.Mux()}
	errGroup.Go(func() error {
		if err := srv.Serve(httpL); err != nil {
			return err
		}
		return nil
	})

	i.logger.Info("listen http on", httpL.Addr().String())
	errGroup.Go(func() error {
		<-ictx.Done()
		return srv.Close()
	})
	return errGroup.Wait()
}

func (i *App) runRTMP(ctx context.Context, rtmpEndpoint string) error {
	errGroup, ictx := errgroup.WithContext(ctx)
	rtmpL, err := net.Listen("tcp", rtmpEndpoint)
	if err != nil {
		return err
	}

	ingressSrv := rtmp.Server{Handler: i.rtmpHandler}
	errGroup.Go(func() error {
		if err := ingressSrv.Serve(rtmpL); err != nil {
			return err
		}
		return nil
	})
	i.logger.Info("listen rtmp-ingress on:", rtmpL.Addr().String())

	errGroup.Go(func() error {
		<-ictx.Done()
		return ingressSrv.Close()
	})
	return errGroup.Wait()
}

///////////////////////////////////////
//
//
//
//
//

func Start(conf Config) {

	// setup context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)

	//setup filesystem
	afs := initFileSystem(conf.Filesystem.Path)

	//setup logger
	logger := initLogger(&conf)

	// kv store
	kv := kvstore.Open()

	//sqlite setup
	dsn := conf.SQL.DSN
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		logger.Fatal(err)
		return
	}

	//GRPC
	grpcAddr := fmt.Sprintf(":%v", conf.Rpc.Port)
	gs := monolith.NewRPC(grpcAddr)
	reflection.Register(gs)

	//HTTP
	httpAddr := fmt.Sprintf(":%v", conf.Http.Port)
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	hs := monolith.NewHTTP(httpAddr, mux)

	//rtmp
	rtmpAddr := fmt.Sprintf(":%v", conf.Rtmp.Port)
	h := rtmp.NewHandler()

	//observer
	obs, err := observer.New(ctx)
	if err != nil {
		logger.Panic(err)
	}
	defer obs.Shutdown(ctx)

	//infrastructure Instance
	infra := App{
		fs:          afs,
		kv:          kv,
		logger:      logger,
		rpc:         gs,
		http:        hs,
		rtmpHandler: h,
		db:          db,

		minioConf: &monolith.MinioConfig{
			Addr:      conf.Minio.Address,
			AccessKey: conf.Minio.AccessKey,
			SecretKey: conf.Minio.SecretKey,
		},

		observer: obs,
	}

	go func() {
		//  shutdown
		s := <-sigC
		cancel()
		infra.logger.Println("got signal", s)
	}()

	// media module
	media.StartModule(&infra)
	/* rtmp module */
	ingress.InitModule(ctx, &infra)
	/* user module */
	users.StartModule(ctx, &infra)
	/* streaming module */
	streaming.StartModule(ctx, &infra)
	/* WebServer (ui) */
	web.StartModule(&infra)

	if err := infra.Run(ctx, grpcAddr, httpAddr, rtmpAddr); err != nil {
		infra.logger.Errorf("shutdown server : %v ", err)
	}

}
