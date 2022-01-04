package api

import (
	"context"
	"net/http"

	coreMiddleware "github.com/facily-tech/go-core/http/server/middleware"

	"github.com/facily-tech/go-scaffold/internal/container"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote/transport"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Handler(ctx context.Context, dep *container.Dependency) http.Handler {
	r := chi.NewMux()

	r.Use(coreMiddleware.Logger(dep.Components.Log))
	r.Use(coreMiddleware.Recoverer(dep.Components.Log))
	r.Use(middleware.RequestID)
	r.Use(dep.Components.Tracer.Middleware)

	r.Handle("/metrics", promhttp.Handler())

	quoteHandler := transport.NewHTTPHandler(dep.Services.Quote)
	r.Mount("/v1/quote", quoteHandler)

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("foo")
	})

	return r
}
