package main

import (
	"context"
	"log"
	"time"

	"github.com/facily-tech/go-core/env"
	apiServer "github.com/facily-tech/go-core/http/server"
	"github.com/facily-tech/go-core/types"
	"github.com/facily-tech/go-scaffold/internal/api"
	"github.com/facily-tech/go-scaffold/internal/config"
	"github.com/facily-tech/go-scaffold/internal/container"

	_ "github.com/golang/mock/mockgen/model"
)

func main() {
	// root context of application
	ctx := context.Background()

	ctx = context.WithValue(ctx, types.ContextKey(types.Version), config.NewVersion())
	ctx = context.WithValue(ctx, types.ContextKey(types.StartedAt), time.Now())

	ctx, dep, err := container.New(ctx)
	if err != nil {
		log.Fatal(err) // log might not be started and because of that dep might not exist
	}

	apiConfig := apiServer.Config{}
	err = env.LoadEnv(ctx, &apiConfig, apiServer.PrefixConfig)
	if err != nil {
		log.Fatal(err)
	}

	apiServer.Run(
		ctx,
		apiConfig,
		api.Handler(ctx, dep),
		dep.Components.Log,
	)

	dep.Components.Tracer.Close()
}
