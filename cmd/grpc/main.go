package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	coreLog "github.com/facily-tech/go-core/log"
	scaffolding "github.com/facily-tech/go-scaffold"
	"github.com/facily-tech/go-scaffold/internal/config"
	"github.com/facily-tech/go-scaffold/internal/container"
	"github.com/facily-tech/go-scaffold/pkg/core/types"
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

	ctx, dep, err := container.New(ctx, scaffolding.Embeds)
	if err != nil {
		log.Fatal(err)
	}

	run(ctx, dep)
}

func run(ctx context.Context, dep *container.Dependency) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer signal.Stop(interrupt)

	server, err := net.Listen("tcp", dep.Components.Viper.GetString("SERVER_BIND"))
	if err != nil {
		dep.Components.Log.Error(
			ctx,
			"Unable to listen",
			coreLog.Any("SERVER_BIND", dep.Components.Viper.GetString("SERVER_BIND")),
			coreLog.Error(err),
		)
		return
	}

	dep.Components.Log.Info(
		ctx,
		"Starting grpc server",
		coreLog.Any("SERVER_BIND", dep.Components.Viper.GetString("SERVER_BIND")),
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

	dep.Components.Trace.Close()
}
