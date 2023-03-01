package quote

import (
	"context"
	"database/sql"

	"sync"

	"github.com/facily-tech/go-core/log"
	"github.com/google/uuid"
)

// RepositoryI is a interface to communicate with a external source of data (ex: Postgres, Firebase FireStore or an API)
// It is using the concept of having a readable interface called "Querier"
// and a Writable interface called "Execer", which exec actions into the external source of data.
type RepositoryI interface {
	// Queries is a "Readeble" interface responsible to read data from source
	Querier

	// Execer is a "Writable" interface responsible for write data into source
	Execer
}

type Querier interface {
	FindByID(context.Context, uuid.UUID) (Quote, error)
}

type Execer interface {
	Upsert(context.Context, *Quote) error
	Delete(context.Context, uuid.UUID) error
}

type RepositoryMemory struct {
	mu  sync.Mutex
	db  map[uuid.UUID]Quote
	log log.Logger
}

func NewRepository(logger log.Logger) *RepositoryMemory {
	return &RepositoryMemory{
		db:  map[uuid.UUID]Quote{},
		log: logger,
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
