package api

import (
	"context"
	"net/http"

	coreMiddleware "github.com/<REPO>/go-core/http/server/middleware"

	"github.com/<REPO>/go-scaffold/internal/container"
	"github.com/<REPO>/go-scaffold/pkg/domains/quote/transport"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Handler(ctx context.Context, dep *container.Dependency) http.Handler {
	r := chi.NewMux()

	r.Use(dep.Components.Tracer.Middleware)             // must be first
	r.Use(middleware.RequestID)                         // must be second
	r.Use(coreMiddleware.Logger(dep.Components.Log))    // must be third
	r.Use(coreMiddleware.Recoverer(dep.Components.Log)) // must be forty

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/health", fun(c *gin.Context) {
		c.JSON(200, gin.H{
		  "status": "UP",
		})

	quoteHandler := transport.NewHTTPHandler(dep.Services.Quote)
	r.Mount("/v1/quote", quoteHandler)

	return r
}
