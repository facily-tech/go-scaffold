package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	scaffolding "github.com/facily-tech/go-scaffold"
	"github.com/facily-tech/go-scaffold/internal/config"
	"github.com/facily-tech/go-scaffold/internal/container"
	"github.com/facily-tech/go-scaffold/pkg/core/types"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote/transport"
	pb "github.com/facily-tech/proto-examples/go-scaffold/build/go/quote"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// root context of application
	ctx := context.Background()

	ctx = context.WithValue(ctx, types.ContextKey(types.Version), config.NewVersion())
	ctx = context.WithValue(ctx, types.ContextKey(types.StartedAt), time.Now())

	ctx, dep, err := container.New(ctx, scaffolding.Embeds)
	if err != nil {
		log.Fatal(err)
	}

	run(ctx, dep)
}

func run(_ context.Context, dep *container.Dependency) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer signal.Stop(interrupt)

	server, err := net.Listen("tcp", viper.GetString("SERVER_BIND"))
	if err != nil {
		dep.Components.Log.Error(
			"Unable to listen",
			zap.String("SERVER_BIND", viper.GetString("SERVER_BIND")),
			zap.Error(err),
		)
		return
	}

	dep.Components.Log.Info(
		"Starting grpc server",
		zap.String("SERVER_BIND", viper.GetString("SERVER_BIND")),
	)

	defer server.Close()
	grpcServer := grpc.NewServer()
	pb.RegisterQuoteServiceServer(grpcServer, transport.NewGRPCServer(dep.Services.Quote))
	reflection.Register(grpcServer)
	go func() {
		if err := grpcServer.Serve(server); err != nil {
			dep.Components.Log.Error("Unexpected grpc server end", zap.Error(err))
		}
	}()

	<-interrupt
	dep.Components.Log.Info("Stopping grpc server")
	grpcServer.GracefulStop()
}
