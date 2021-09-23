package quote

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func createCounter(t *testing.T) (*ServiceI_Mock, prometheus.Counter) {
	service := NewServiceI_Mock(t)
	registry := prometheus.NewRegistry()
	return service, promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Name: UpsertMetricName,
			Help: UpsertHelp,
		},
	)
}

func TestNewCountMetric(t *testing.T) {
	service, counter := createCounter(t)

	type args struct {
		quoteService ServiceI
		counter      prometheus.Counter
	}
	tests := []struct {
		name string
		args args
		want *CountMetric
	}{
		{
			name: "success, create new upsert count metric",
			args: args{quoteService: service, counter: counter},
			want: &CountMetric{
				count: counter,
				next:  service,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCountMetric(tt.args.quoteService, tt.args.counter)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCountMetric_FindByID(t *testing.T) {
	service, counter := createCounter(t)

	type fields struct {
		count prometheus.Counter
		next  ServiceI
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Quote
		init   func(ctx context.Context, id uuid.UUID)
		err    error
	}{
		{
			name:   "success, call find",
			fields: fields{count: counter, next: service},
			args:   args{ctx: context.Background(), id: testQuote.ID},
			init: func(ctx context.Context, id uuid.UUID) {
				service.On_FindByID().Args(ctx, id).Rets(testQuote, nil).Times(1)
			},
			want: testQuote,
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := &CountMetric{
				count: tt.fields.count,
				next:  tt.fields.next,
			}
			tt.init(tt.args.ctx, tt.args.id)
			got, err := mw.FindByID(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, service.Mock.AssertExpectations())
		})
	}
}

func TestCountMetric_Upsert(t *testing.T) {
	service, counter := createCounter(t)

	type fields struct {
		count prometheus.Counter
		next  ServiceI
	}
	type args struct {
		ctx context.Context
		q   *Quote
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		init   func(ctx context.Context, q *Quote)
		err    error
	}{
		{
			name:   "success, call upsert and increase prometheus counter",
			fields: fields{count: counter, next: service},
			args:   args{ctx: context.Background(), q: &testQuote},
			init: func(ctx context.Context, q *Quote) {
				service.On_Upsert().Args(ctx, q).Rets(nil).Times(1)
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := &CountMetric{
				count: tt.fields.count,
				next:  tt.fields.next,
			}
			tt.init(tt.args.ctx, tt.args.q)
			err := mw.Upsert(tt.args.ctx, tt.args.q)
			assert.ErrorIs(t, err, tt.err)
			assert.NoError(t, service.Mock.AssertExpectations())
			assert.Equal(t, 1, testutil.CollectAndCount(counter, UpsertMetricName))
		})
	}
}

func TestCountMetric_Delete(t *testing.T) {
	service, counter := createCounter(t)
	type fields struct {
		count prometheus.Counter
		next  ServiceI
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		init   func(context.Context, uuid.UUID)
		err    error
	}{
		{
			name:   "success, call Delete",
			fields: fields{count: counter, next: service},
			args:   args{ctx: context.Background(), id: testQuote.ID},
			init: func(c context.Context, u uuid.UUID) {
				service.On_Delete().Args(c, u).Rets(nil).Times(1)
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := &CountMetric{
				count: tt.fields.count,
				next:  tt.fields.next,
			}
			tt.init(tt.args.ctx, tt.args.id)
			err := mw.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.err)
			assert.NoError(t, service.Mock.AssertExpectations())
		})
	}
}
