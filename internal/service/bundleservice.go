package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

type bundleService struct {
	BR domain.BundleRepository
}

func NewBundleService(br domain.BundleRepository) *bundleService {
	return &bundleService{
		BR: br,
	}
}

func (bs bundleService) FetchByID(ctx context.Context, bid int64) (*domain.Bundle, error) {
	return nil, nil
}

func (bs bundleService) FetchAll(ctx context.Context) (*[]domain.Bundle, error) {
	return nil, nil
}

func (bs bundleService) Store(ctx context.Context, bundle *domain.Bundle) error {
	return nil
}

func (bs bundleService) Remove(ctx context.Context, bid int64) error {
	return nil
}

func (bs bundleService) Update(ctx context.Context, bundle *domain.Bundle) error {
	return nil
}
