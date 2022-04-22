package quote

import (
	"context"
	"database/sql"
	"testing"

	"github.com/facily-tech/go-core/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testQuote = Quote{ID: uuid.New(), Content: "IRC is just multiplayer notepad."}

func TestNewService(t *testing.T) {
	l, err := log.NewLoggerZap(log.ZapConfig{})
	assert.NoError(t, err)

	inmem := NewRepository(l)
	type args struct {
		repository RepositoryI
		log        log.Logger
	}
	tests := []struct {
		name string
		args args
		want *Service
		err  error
	}{
		{
			name: "success, repository is initiliazed",
			args: args{
				repository: inmem,
				log:        l,
			},
			want: &Service{
				repository: inmem,
				log:        l,
			},
			err: nil,
		},
		{
			name: "fail, repository is empty",
			args: args{
				repository: nil,
				log:        nil,
			},
			want: nil,
			err:  ErrEmptyRepository,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(tt.args.repository, tt.args.log)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_FindByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	log, _ := log.NewLoggerZap(log.ZapConfig{})
	t.Run("success, create and find the quote", func(t *testing.T) {
		inmem := NewRepository(log)

		tt := struct {
			args args
			want Quote
			err  error
		}{
			args: args{ctx: context.Background(), id: testQuote.ID},
			want: testQuote,
			err:  nil,
		}

		s := &Service{
			repository: inmem,
		}
		err := s.Upsert(tt.args.ctx, &tt.want)
		assert.NoError(t, err)

		got, err := s.FindByID(tt.args.ctx, tt.args.id)
		assert.ErrorIs(t, err, tt.err)
		assert.Equal(t, tt.want, got)
	})

	t.Run("fail, try to find the quote which is not created before", func(t *testing.T) {
		inmem := NewRepository(log)

		tt := struct {
			args args
			want Quote
			err  error
		}{
			args: args{ctx: context.Background(), id: testQuote.ID},
			want: Quote{},
			err:  sql.ErrNoRows,
		}

		s := &Service{
			repository: inmem,
		}

		got, err := s.FindByID(tt.args.ctx, tt.args.id)
		assert.ErrorIs(t, err, tt.err)
		assert.Equal(t, tt.want, got)
	})
}

func TestService_Upsert(t *testing.T) {
	log, _ := log.NewLoggerZap(log.ZapConfig{})
	inmem := NewRepository(log)

	type fields struct {
		repository RepositoryI
	}
	type args struct {
		ctx context.Context
		q   *Quote
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		init   func(*testing.T, ServiceI)
		after  func(*testing.T, ServiceI, uuid.UUID, Quote)
		err    error
	}{
		{
			name:   "success, insert new quote",
			fields: fields{repository: inmem},
			args:   args{ctx: context.Background(), q: &testQuote},
			init:   func(t *testing.T, s ServiceI) {},
			after:  func(t *testing.T, si ServiceI, id uuid.UUID, q Quote) {},
			err:    nil,
		},
		{
			name:   "success, insert new quote",
			fields: fields{repository: inmem},
			args:   args{ctx: context.Background(), q: &Quote{ID: testQuote.ID, Content: "changed"}},
			init: func(t *testing.T, s ServiceI) {
				err := s.Upsert(context.Background(), &testQuote)
				assert.NoError(t, err)
			},
			after: func(t *testing.T, s ServiceI, id uuid.UUID, quote Quote) {
				q, err := s.FindByID(context.Background(), id)
				assert.NoError(t, err)
				assert.Equal(t, quote, q)
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repository: tt.fields.repository,
			}
			tt.init(t, s)
			err := s.Upsert(tt.args.ctx, tt.args.q)
			assert.ErrorIs(t, err, tt.err)
			tt.after(t, s, tt.args.q.ID, *tt.args.q)
		})
	}
}

func TestService_Delete(t *testing.T) {
	log, _ := log.NewLoggerZap(log.ZapConfig{})
	inmem := NewRepository(log)
	type fields struct {
		repository RepositoryI
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		after  func(ServiceI, uuid.UUID)
		err    error
	}{
		{
			name:   "success, create and remove quote",
			fields: fields{repository: inmem},
			args:   args{ctx: context.Background(), id: testQuote.ID},
			after:  func(s ServiceI, id uuid.UUID) {},
			err:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repository: tt.fields.repository,
			}
			err := s.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.err)
			tt.after(s, tt.args.id)
		})
	}
}
