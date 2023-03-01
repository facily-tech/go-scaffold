package transport

import (
	"context"
	"encoding/json"
	"log"
	stdHTTP "net/http"

	"github.com/facily-tech/go-scaffold/pkg/domains/quote"
	"github.com/go-chi/chi/v5"
	"github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
)

func NewHTTPHandler(svc quote.ServiceI) stdHTTP.Handler {
	options := []http.ServerOption{
		http.ServerErrorEncoder(errorHandler),
	}

	find := http.NewServer(
		quote.FindByID(svc),
		decodeFindRequest,
		codeHTTP{200}.encodeResponse,
		options...,
	)

	upsert := http.NewServer(
		quote.Upsert(svc),
		decodeUpsertRequest,
		codeHTTP{201}.encodeResponse,
		options...,
	)

	delete := http.NewServer(
		quote.DeleteByID(svc),
		decodeDeleteRequest,
		func(ctx context.Context, w stdHTTP.ResponseWriter, response interface{}) error {
			return nil
		},
		options...,
	)

	r := chi.NewRouter()

	r.Get("/{id}", find.ServeHTTP)
	r.Post("/", upsert.ServeHTTP)
	r.Delete("/{id}", delete.ServeHTTP)

	return r
}

func decodeUpsertRequest(_ context.Context, r *stdHTTP.Request) (interface{}, error) {
	req := quote.UpsertRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeFindRequest(_ context.Context, r *stdHTTP.Request) (interface{}, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return nil, err
	}

	return quote.FindByIDRequest{ID: id}, nil
}

func decodeDeleteRequest(_ context.Context, r *stdHTTP.Request) (interface{}, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return nil, err
	}

	return quote.DeleteRequest{ID: id}, nil
}

// codeHTTP encodes the default HTTP code. if you have/need to handle more status
// codes create your own EncodeResponseFunc.
type codeHTTP struct {
	int
}

func (c codeHTTP) encodeResponse(_ context.Context, w stdHTTP.ResponseWriter, input interface{}) error {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(c.int)
	return json.NewEncoder(w).Encode(input)
}

func errorHandler(_ context.Context, err error, w stdHTTP.ResponseWriter) {
	resp, code := quote.RESTErrorBussines.ErrorProcess(err)

	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": resp}); err != nil {
		log.Printf("Encoding error, nothing much we can do: %v", err)
	}
}
