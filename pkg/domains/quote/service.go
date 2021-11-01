package quote

import (
	"context"

	"github.com/facily-tech/go-scaffold/pkg/domains/quote/models"
	"github.com/google/uuid"
)

//go:generate genmock -search.name=ServiceI -print.place.in_package -print.file.test
type ServiceI interface {
	FindByID(context.Context, uuid.UUID) (models.Quote, error)
	Upsert(context.Context, *models.Quote) error
	Delete(context.Context, uuid.UUID) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) (*Service, error) {
	if repository == nil {
		return nil, ErrEmptyRepository
	}
	return &Service{repository: repository}, nil
}

func (s *Service) FindByID(ctx context.Context, id uuid.UUID) (models.Quote, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Upsert(ctx context.Context, q *models.Quote) error {
	return s.repository.Upsert(ctx, q)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
