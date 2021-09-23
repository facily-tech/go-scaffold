package quote

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
)

type RepositoryI interface {
	FindByID(context.Context, uuid.UUID) (Quote, error)
	Upsert(context.Context, *Quote) error
	Delete(context.Context, uuid.UUID) error
}

type RepositoryMemory struct {
	mu sync.Mutex
	db map[uuid.UUID]Quote
}

func NewRepository() *RepositoryMemory {
	return &RepositoryMemory{
		db: map[uuid.UUID]Quote{},
	}
}

func (r *RepositoryMemory) FindByID(ctx context.Context, id uuid.UUID) (Quote, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	quote, exist := r.db[id]
	if !exist {
		return Quote{}, sql.ErrNoRows
	}

	return quote, nil
}

func (r *RepositoryMemory) Upsert(ctx context.Context, q *Quote) error {
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
