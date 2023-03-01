package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/facily-tech/go-core/env"
	coreLog "github.com/facily-tech/go-core/log"
	"github.com/facily-tech/go-core/types"
	"github.com/facily-tech/go-scaffold/internal/config"
	"github.com/facily-tech/go-scaffold/internal/container"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote/transport"
	pb "github.com/facily-tech/proto-examples/go-scaffold/build/go/quote"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/golang/mock/mockgen/model"
)

func main() {
	// root context of application
	ctx := context.Background()

	ctx = context.WithValue(ctx, types.ContextKey(types.Version), config.NewVersion())
	ctx = context.WithValue(ctx, types.ContextKey(types.StartedAt), time.Now())

	ctx, dep, err := container.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	run(ctx, dep)
}

var grpcPrefixConfig = "GRPC_"

type grpcConfig struct {
	ServerBind string `env:"SERVER_BIND"`
}

func run(ctx context.Context, dep *container.Dependency) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer signal.Stop(interrupt)

	grpcConfig := grpcConfig{}
	err := env.LoadEnv(ctx, &grpcConfig, grpcPrefixConfig)

	if err != nil {
		dep.Components.Log.Error(
			ctx,
			"Unable to load config",
			coreLog.Error(err),
		)
		return
	}

	server, err := net.Listen("tcp", grpcConfig.ServerBind)
	if err != nil {
		dep.Components.Log.Error(
			ctx,
			"Unable to listen",
			coreLog.Any("SERVER_BIND", grpcConfig.ServerBind),
			coreLog.Error(err),
		)
		return
	}

	dep.Components.Log.Info(
		ctx,
		"Starting grpc server",
		coreLog.Any("SERVER_BIND", grpcConfig.ServerBind),
	)

	defer server.Close()
	grpcServer := grpc.NewServer()
	pb.RegisterQuoteServiceServer(grpcServer, transport.NewGRPCServer(dep.Services.Quote))
	reflection.Register(grpcServer)
	go func() {
		if err := grpcServer.Serve(server); err != nil {
			dep.Components.Log.Error(ctx, "Unexpected grpc server end", coreLog.Error(err))
		}
	}()

	<-interrupt
	dep.Components.Log.Info(ctx, "Stopping grpc server")
	grpcServer.GracefulStop()

	dep.Components.Tracer.Close()
}
