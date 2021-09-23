package api

import (
	"context"
	"net/http"

	"github.com/facily-tech/go-scaffold/internal/container"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote/transport"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Handler(ctx context.Context, srv *container.Services) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Handle("/metrics", promhttp.Handler())

	quoteHandler := transport.NewHTTPHandler(srv.Quote)
	r.Mount("/v1/quote", quoteHandler)

	return r
}
