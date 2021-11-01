package mem

import (
	"context"
	"database/sql"
	"sync"

	"github.com/facily-tech/go-scaffold/pkg/domains/quote/models"
	"github.com/google/uuid"
)

type RepositoryMemory struct {
	mu sync.Mutex
	db map[uuid.UUID]models.Quote
}

func NewRepository() *RepositoryMemory {
	return &RepositoryMemory{
		db: map[uuid.UUID]models.Quote{},
	}
}

func (r *RepositoryMemory) FindByID(ctx context.Context, id uuid.UUID) (models.Quote, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	quote, exist := r.db[id]
	if !exist {
		return quote, sql.ErrNoRows
	}

	return quote, nil
}

func (r *RepositoryMemory) Upsert(ctx context.Context, q *models.Quote) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db[q.ID] = *q

	return nil
}

func (r *RepositoryMemory) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.db, id)

	return nil
}
