package quote

import (
	"context"
	"database/sql"
	"testing"

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
				s.On_FindByID().Args(ctx, testQuote.ID).Rets(testQuote, nil)
			},
			args: args{
				svc:     NewServiceI_Mock(t),
				request: FindByIDRequest{ID: testQuote.ID},
				ctx:     context.Background(),
			},
			response: JSONResponse{
				ID:      &testQuote.ID,
				Content: testQuote.Content,
			},
		},
		{
			name: "fail, returning error",
			init: func(s *ServiceI_Mock, ctx context.Context) {
				s.On_FindByID().Args(ctx, testQuote.ID).Rets(Quote{}, sql.ErrNoRows)
			},
			err: sql.ErrNoRows,
			args: args{
				svc:     NewServiceI_Mock(t),
				request: FindByIDRequest{ID: testQuote.ID},
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
				s.On_Upsert().Args(ctx, &testQuote).Rets(nil)
			},
			args: args{
				svc:     NewServiceI_Mock(t),
				request: UpsertRequest{JSONRequest: JSONRequest{ID: &testQuote.ID, Content: testQuote.Content}},
				ctx:     context.Background(),
			},
			response: JSONResponse{
				ID:      &testQuote.ID,
				Content: testQuote.Content,
			},
		},
		{
			name: "fail, empty quote content",
			init: func(s *ServiceI_Mock, ctx context.Context) {
				s.On_Upsert().Args(ctx, &Quote{ID: testQuote.ID}).Rets(sql.ErrNoRows)
			},
			err: ErrNew,
			args: args{
				svc:     NewServiceI_Mock(t),
				request: UpsertRequest{JSONRequest: JSONRequest{ID: &testQuote.ID, Content: ""}},
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
				s.On_Delete().Args(ctx, testQuote.ID).Rets(nil)
			},
			args: args{
				svc:     NewServiceI_Mock(t),
				request: DeleteRequest{ID: testQuote.ID},
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
