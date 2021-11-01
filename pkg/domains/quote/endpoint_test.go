package quote

import (
	"context"
	"database/sql"
	"testing"

	"github.com/facily-tech/go-scaffold/pkg/domains/quote/models"
	"github.com/stretchr/testify/assert"
)

func TestFindByID(t *testing.T) {
	type args struct {
		svc     *ServiceI_Mock
		request interface{}
		ctx     context.Context
	}
	tests := []struct {
		name     string
		args     args
		response interface{}
		init     func(*ServiceI_Mock, context.Context)
		err      error
	}{
		{
			name: "success",
			init: func(s *ServiceI_Mock, ctx context.Context) {
				s.On_FindByID().Args(ctx, models.TestQuote.ID).Rets(models.TestQuote, nil)
			},
			args: args{
				svc:     NewServiceI_Mock(t),
				request: FindByIDRequest{ID: models.TestQuote.ID},
				ctx:     context.Background(),
			},
			response: models.JSONResponse{
				ID:      &models.TestQuote.ID,
				Content: models.TestQuote.Content,
			},
		},
		{
			name: "fail, returning error",
			init: func(s *ServiceI_Mock, ctx context.Context) {
				s.On_FindByID().Args(ctx, models.TestQuote.ID).Rets(models.Quote{}, sql.ErrNoRows)
			},
			err: sql.ErrNoRows,
			args: args{
				svc:     NewServiceI_Mock(t),
				request: FindByIDRequest{ID: models.TestQuote.ID},
				ctx:     context.Background(),
			},
			response: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.init(tt.args.svc, tt.args.ctx)
			response, err := FindByID(tt.args.svc)(tt.args.ctx, tt.args.request)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.response, response)
		})
	}
}

func TestUpsert(t *testing.T) {
	type args struct {
		svc     *ServiceI_Mock
		request interface{}
		ctx     context.Context
	}
	tests := []struct {
		name     string
		args     args
		response interface{}
		init     func(*ServiceI_Mock, context.Context)
		err      error
	}{
		{
			name: "success",
			init: func(s *ServiceI_Mock, ctx context.Context) {
				s.On_Upsert().Args(ctx, &models.TestQuote).Rets(nil)
			},
			args: args{
				svc:     NewServiceI_Mock(t),
				request: UpsertRequest{JSONRequest: JSONRequest{ID: &models.TestQuote.ID, Content: models.TestQuote.Content}},
				ctx:     context.Background(),
			},
			response: models.JSONResponse{
				ID:      &models.TestQuote.ID,
				Content: models.TestQuote.Content,
			},
		},
		{
			name: "fail, empty quote content",
			init: func(s *ServiceI_Mock, ctx context.Context) {
				s.On_Upsert().Args(ctx, &models.Quote{ID: models.TestQuote.ID}).Rets(sql.ErrNoRows)
			},
			err: models.ErrNew,
			args: args{
				svc:     NewServiceI_Mock(t),
				request: UpsertRequest{JSONRequest: JSONRequest{ID: &models.TestQuote.ID, Content: ""}},
				ctx:     context.Background(),
			},
			response: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.init(tt.args.svc, tt.args.ctx)
			response, err := Upsert(tt.args.svc)(tt.args.ctx, tt.args.request)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.response, response)
		})
	}
}

func TestDeleteByID(t *testing.T) {
	type args struct {
		svc     *ServiceI_Mock
		request interface{}
		ctx     context.Context
	}
	tests := []struct {
		name     string
		args     args
		response interface{}
		init     func(*ServiceI_Mock, context.Context)
		err      error
	}{
		{
			name: "success",
			init: func(s *ServiceI_Mock, ctx context.Context) {
				s.On_Delete().Args(ctx, models.TestQuote.ID).Rets(nil)
			},
			args: args{
				svc:     NewServiceI_Mock(t),
				request: DeleteRequest{ID: models.TestQuote.ID},
				ctx:     context.Background(),
			},
			response: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.init(tt.args.svc, tt.args.ctx)
			response, err := DeleteByID(tt.args.svc)(tt.args.ctx, tt.args.request)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.response, response)
		})
	}
}
