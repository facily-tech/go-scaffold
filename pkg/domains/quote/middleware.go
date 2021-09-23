package quote

import (
	"context"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	UpsertMetricName string = "upsert_total"
	UpsertHelp       string = "total number of calls to upsert"
)

// CountMetric is a domain level metric middleware using prometheus.
type CountMetric struct {
	count prometheus.Counter
	next  ServiceI
}

func NewCountMetric(quoteService ServiceI, counter prometheus.Counter) *CountMetric {
	return &CountMetric{
		count: counter,
		next:  quoteService,
	}
}

func (mw *CountMetric) FindByID(ctx context.Context, id uuid.UUID) (Quote, error) {
	return mw.next.FindByID(ctx, id)
}

func (mw *CountMetric) Upsert(ctx context.Context, q *Quote) error {
	defer mw.count.Inc()
	return mw.next.Upsert(ctx, q)
}

func (mw *CountMetric) Delete(ctx context.Context, id uuid.UUID) error {
	return mw.next.Delete(ctx, id)
}
