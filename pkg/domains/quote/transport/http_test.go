package transport

import (
	"context"
	"errors"
	"io"
	stdHTTP "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/facily-tech/go-scaffold/pkg/domains/quote"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_decodeUpsertRequest(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *stdHTTP.Request
	}
	tests := []struct {
		name string
		args args
		want interface{}
		err  error
	}{
		{
			name: "success, decode request for upsert",
			args: args{
				ctx: context.Background(),
				r: httptest.NewRequest(
					"GET",
					"/",
					strings.NewReader(`{"content": "Quote"}`),
				),
			},
			want: quote.UpsertRequest{JSONRequest: quote.JSONRequest{Content: "Quote"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeUpsertRequest(tt.args.ctx, tt.args.r)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_decodeFindRequest(t *testing.T) {
	id := uuid.New()
	type args struct {
		ctx context.Context
		r   *stdHTTP.Request
	}
	tests := []struct {
		name string
		args args
		want interface{}
		init func(r *stdHTTP.Request) *stdHTTP.Request
		err  error
	}{
		{
			name: "success, decode request for find",
			args: args{
				ctx: context.Background(),
				r: httptest.NewRequest(
					"GET",
					"/"+id.String(),
					strings.NewReader(`{"content": "Quote"}`),
				),
			},
			init: func(r *stdHTTP.Request) *stdHTTP.Request {
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", id.String())
				return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
			},
			want: quote.FindByIDRequest{ID: id},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.init(tt.args.r)
			got, err := decodeFindRequest(tt.args.ctx, r)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func Test_decodeDeleteRequest(t *testing.T) {
	id := uuid.New()
	type args struct {
		ctx context.Context
		r   *stdHTTP.Request
	}
	tests := []struct {
		name string
		args args
		want interface{}
		init func(r *stdHTTP.Request) *stdHTTP.Request
		err  error
	}{
		{
			name: "success, decode request for delete",
			args: args{
				ctx: context.Background(),
				r: httptest.NewRequest(
					"GET",
					"/"+id.String(),
					nil,
				),
			},
			init: func(r *stdHTTP.Request) *stdHTTP.Request {
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", id.String())
				return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
			},
			want: quote.DeleteRequest{ID: id},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.init(tt.args.r)
			got, err := decodeDeleteRequest(tt.args.ctx, r)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func Test_codeHTTP_encodeResponse(t *testing.T) {
	type fields struct {
		int int
	}
	type args struct {
		ctx   context.Context
		w     *httptest.ResponseRecorder
		input interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantCode   int
		wantHeader stdHTTP.Header
		err        error
	}{
		{
			name:   "success",
			fields: fields{int: 200},
			args: args{
				ctx:   context.Background(),
				w:     httptest.NewRecorder(),
				input: map[string]string{"key": "value"},
			},
			wantCode:   200,
			wantHeader: stdHTTP.Header{"Content-Type": []string{"application/json; charset=UTF-8"}},
			want:       `{"key":"value"}` + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := codeHTTP{
				int: tt.fields.int,
			}
			err := c.encodeResponse(tt.args.ctx, tt.args.w, tt.args.input)
			assert.ErrorIs(t, err, tt.err)
			result, err := io.ReadAll(tt.args.w.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(result))
			assert.Equal(t, tt.wantCode, tt.args.w.Code)
			assert.Equal(t, len(tt.args.w.Header()), len(tt.args.w.Header()))
			for k, v := range tt.args.w.Header() {
				assert.Equal(t, tt.wantHeader[k], v, k)
			}
		})
	}
}

func Test_errorHandler(t *testing.T) {
	type args struct {
		ctx context.Context
		err error
		w   *httptest.ResponseRecorder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				err: errors.New("random error"), //nolint:goerr113 //allowed for testing
				w:   httptest.NewRecorder(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorHandler(tt.args.ctx, tt.args.err, tt.args.w)
		})
	}
}
