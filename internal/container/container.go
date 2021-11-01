package container

import (
	"context"
	"embed"
	"os"

	"github.com/facily-tech/go-scaffold/pkg/core/log"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote/repository/sql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

// Components are a like service, but it doesn't include business case
// Or domains, but likely used by multiple domains
type components struct {
	Log *zap.Logger
	// Include your new components bellow
}

// Services hold the business case, and make the bridge between
// Controllers and Domains
type Services struct {
	Quote quote.ServiceI
	// Include your new services bellow
}

type Dependency struct {
	Components components
	Services   Services
}

func New(ctx context.Context, embs embed.FS) (context.Context, *Dependency, error) {
	cmp, err := setupComponents(ctx, embs)
	if err != nil {
		return nil, nil, err
	}

	// TODO sql.MakeConfig
	repository, err := sql.NewAndMigrate(ctx, os.Getenv("DB_DSN"))
	if err != nil {
		return nil, nil, err
	}
	quoteService, err := quote.NewService(repository)
	if err != nil {
		return nil, nil, err
	}

	srv := Services{
		quote.NewCountMetric(
			quoteService,
			promauto.NewCounter(prometheus.CounterOpts{
				Name: quote.UpsertMetricName,
				Help: quote.UpsertHelp,
			}),
		),
		// include services initialized above here
	}

	dep := Dependency{
		Components: *cmp,
		Services:   srv,
	}

	return ctx, &dep, err
}

func setupComponents(_ context.Context, _ embed.FS) (*components, error) {
	log, err := log.NewLogger(true)
	if err != nil {
		return nil, err
	}

	return &components{
		// include components initialized above here
		Log: log,
	}, nil
}
