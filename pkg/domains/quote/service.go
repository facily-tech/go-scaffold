package quote

import (
	"context"

	"github.com/facily-tech/go-core/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	quote, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return quote, errors.Wrap(err, "service failed to find by ID")
	}

	return quote, nil
}

func (s *Service) Upsert(ctx context.Context, q *Quote) error {
	if err := s.repository.Upsert(ctx, q); err != nil {
		return errors.Wrap(err, "service failed to upsert")
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return errors.Wrap(err, "service failed to delete")
	}

	return nil
}
