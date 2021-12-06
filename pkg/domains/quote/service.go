package quote

import (
	"context"

	"github.com/facily-tech/go-core/log"
	"github.com/google/uuid"
)

//go:generate mockgen -destination service_mock.go -package=quote -source=service.go
type ServiceI interface {
	FindByID(context.Context, uuid.UUID) (Quote, error)
	Upsert(context.Context, *Quote) error
	Delete(context.Context, uuid.UUID) error
}

type Service struct {
	repository RepositoryI
	log        log.Logger
}

func NewService(repository RepositoryI, log log.Logger) (*Service, error) {
	if repository == nil {
		return nil, ErrEmptyRepository
	}
	return &Service{
		repository: repository,
		log:        log,
	}, nil
}

func (s *Service) FindByID(ctx context.Context, id uuid.UUID) (Quote, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Upsert(ctx context.Context, q *Quote) error {
	return s.repository.Upsert(ctx, q)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
