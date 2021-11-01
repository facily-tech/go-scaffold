package quote

import (
	"context"

	"github.com/facily-tech/go-scaffold/pkg/domains/quote/models"
	"github.com/google/uuid"
)

type Repository interface {
	Querier
	Execer
}

type Querier interface {
	FindByID(context.Context, uuid.UUID) (models.Quote, error)
}

type Execer interface {
	Upsert(context.Context, *models.Quote) error
	Delete(context.Context, uuid.UUID) error
}
