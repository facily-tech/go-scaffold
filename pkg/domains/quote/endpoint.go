package quote

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
		req, ok := request.(FindByIDRequest)
		if !ok {
			return nil, errors.Wrap(ErrTypeAssertion, "cannot convert request->FindByIDRequest")
		}

		q, err := svc.FindByID(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return q.toJSONResponse(), nil
	}
}

func Upsert(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpsertRequest)
		if !ok {
			return nil, errors.Wrap(ErrTypeAssertion, "cannot convert request->UpsertRequest")
		}

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
		req, ok := request.(DeleteRequest)
		if !ok {
			return nil, errors.Wrap(ErrTypeAssertion, "cannot convert request->Delete")
		}

		// nothing doing anything with error for now
		_ = svc.Delete(ctx, req.ID)
		return nil, nil
	}
}
