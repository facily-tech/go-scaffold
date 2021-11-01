package main

import (
	"context"
	"log"
	"os"
	"time"

	scaffolding "github.com/facily-tech/go-scaffold"
	"github.com/facily-tech/go-scaffold/internal/api"
	"github.com/facily-tech/go-scaffold/internal/config"
	"github.com/facily-tech/go-scaffold/internal/container"
	apiServer "github.com/facily-tech/go-scaffold/pkg/core/http/server"
	"github.com/facily-tech/go-scaffold/pkg/core/types"
	"github.com/spf13/viper"
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

	apiServer.Run(
		ctx,
		apiServer.Config{
			Addr:             os.Getenv("API_HOST_PORT"),
			GracefulDuration: viper.GetDuration("API_GRACEFUL_WAIT_TIME"),
		},
		api.Handler(ctx, &dep.Services),
		dep.Components.Log,
	)
}
