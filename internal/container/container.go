package container

import (
	"context"
	"embed"

	"github.com/facily-tech/go-core/log"
	"github.com/facily-tech/go-core/telemetry"
	"github.com/facily-tech/go-core/types"
	"github.com/facily-tech/go-scaffold/internal/config"
	"github.com/facily-tech/go-scaffold/pkg/core/env"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
)

// Components are a like service, but it doesn't include business case
// Or domains, but likely used by multiple domains
type components struct {
	Viper  *viper.Viper
	Log    log.Logger
	Tracer telemetry.Tracer
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

	quoteService, err := quote.NewService(
		quote.NewRepository(cmp.Log),
		cmp.Log,
	)

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

func setupComponents(ctx context.Context, embedFS embed.FS) (*components, error) {
	version, ok := ctx.Value(types.ContextKey(types.Version)).(*config.Version)
	if !ok {
		return nil, config.ErrVersionTypeAssertion
	}

	vip, err := env.ViperConfig(embedFS)
	if err != nil {
		return nil, err
	}

	tracer := telemetry.NewDataDog(
		telemetry.DataDogConfig{
			Env:     viper.GetString("DD_ENV"),
			Service: viper.GetString("DD_SERVICE"),
			Version: version.GitCommitHash,
		},
	)

	l, err := log.NewLoggerZap(log.ZapConfig{
		Version:           version.GitCommitHash,
		DisableStackTrace: true,
		Tracer:            tracer,
	})

	if err != nil {
		return nil, err
	}

	return &components{
		Viper:  vip,
		Log:    l,
		Tracer: tracer,
		// include components initialized bellow here
	}, nil
}
