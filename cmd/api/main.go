package main

import (
	"context"
	"log"
	"time"

	scaffolding "github.com/facily-tech/go-scaffold"
	"github.com/facily-tech/go-scaffold/internal/api"
	"github.com/facily-tech/go-scaffold/internal/config"
	"github.com/facily-tech/go-scaffold/internal/container"
	apiServer "github.com/facily-tech/go-scaffold/pkg/core/http/server"
	"github.com/facily-tech/go-scaffold/pkg/core/types"

	_ "github.com/golang/mock/mockgen/model"
)

func main() {
	// root context of application
	ctx := context.Background()

	ctx = context.WithValue(ctx, types.ContextKey(types.Version), config.NewVersion())
	ctx = context.WithValue(ctx, types.ContextKey(types.StartedAt), time.Now())

	ctx, dep, err := container.New(ctx, scaffolding.Embeds)
	if err != nil {
		log.Fatal(err) // log might not be started and because of that dep might not exist
	}

	apiServer.Run(
		ctx,
		apiServer.Config{
			Addr:             dep.Components.Viper.GetString("API_HOST_PORT"),
			GracefulDuration: dep.Components.Viper.GetDuration("API_GRACEFUL_WAIT_TIME"),
		},
		api.Handler(ctx, dep),
		dep.Components.Log,
	)

	dep.Components.Tracer.Close()
}
