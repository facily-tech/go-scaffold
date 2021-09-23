package container

import (
	"context"
	"embed"

	"github.com/facily-tech/go-scaffold/pkg/core/env"
	"github.com/facily-tech/go-scaffold/pkg/core/log"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Components are a like service, but it doesn't include business case
// Or domains, but likely used by multiple domains
type components struct {
	Viper *viper.Viper
	Log   *zap.Logger
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

	quoteService, err := quote.NewService(quote.NewRepository())
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

func setupComponents(_ context.Context, embedFS embed.FS) (*components, error) {
	vip, err := env.ViperConfig(embedFS)
	if err != nil {
		return nil, err
	}

	log, err := log.NewLogger(true)
	if err != nil {
		return nil, err
	}

	return &components{
		// include components initialized above here
		Viper: vip,
		Log:   log,
	}, nil
}
