package sql

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote/models"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote/repository/sql/ent"
	"github.com/facily-tech/go-scaffold/pkg/domains/quote/repository/sql/ent/quote"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	_ "github.com/lib/pq" // dialect for postgres
)

const drv = "postgres"

type DB struct {
	client *ent.Client
}

func NewAndMigrate(ctx context.Context, dsn string, opts ...ent.Option) (*DB, error) {
	c, err := ent.Open(drv, dsn, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "can't create DB instance")
	}

	if err := c.Schema.Create(ctx); err != nil {
		return nil, errors.Wrap(err, "can't migrate")
	}

	return &DB{client: c}, nil
}

func entToModel(s ent.Quote) models.Quote {
	return models.Quote{
		ID:      s.ID,
		Content: s.Content,
	}
}

func (db *DB) FindByID(ctx context.Context, id uuid.UUID) (models.Quote, error) {
	dbQuote, err := db.client.Quote.Query().Where(quote.ID(id)).Only(ctx)
	if err != nil {
		return models.Quote{}, errors.Wrap(err, "cannot find only one Quote")
	}

	return entToModel(*dbQuote), nil
}

func (db *DB) Upsert(ctx context.Context, q *models.Quote) error {
	return db.client.
		Quote.
		Create().
		SetID(q.ID).
		SetContent(q.Content).
		OnConflict(
			sql.ConflictColumns("id"),
			sql.ResolveWithNewValues(),
		).
		UpdateNewValues().
		Exec(ctx)
}

func (db *DB) Delete(ctx context.Context, id uuid.UUID) error {
	return db.client.Quote.DeleteOneID(id).Exec(ctx)
}
