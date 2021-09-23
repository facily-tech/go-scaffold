package quote

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type FindByIDRequest struct {
	ID uuid.UUID
}

type UpsertRequest struct {
	JSONRequest
}

type DeleteRequest struct {
	ID uuid.UUID
}

type JSONRequest struct {
	ID      *uuid.UUID `json:"id,omitempty"`
	Content string     `json:"content"`
}

type JSONResponse struct {
	ID      *uuid.UUID `json:"id,omitempty"`
	Content string     `json:"content,omitempty"`
	Error   string     `json:"error,omitempty"`
}

func FindByID(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindByIDRequest)
		q, err := svc.FindByID(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return q.toJSONResponse(), nil
	}
}

func Upsert(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpsertRequest)

		q, err := NewQuote(req.ID, req.Content)
		if err != nil {
			return nil, err
		}

		if err := svc.Upsert(ctx, &q); err != nil {
			return nil, err
		}

		return q.toJSONResponse(), nil
	}
}

func DeleteByID(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)

		// nothing doing anything with error for now
		_ = svc.Delete(ctx, req.ID)
		return nil, nil
	}
}
