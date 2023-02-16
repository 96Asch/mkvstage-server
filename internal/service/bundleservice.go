package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

type bundleService struct {
	br domain.BundleRepository
}

func NewBundleService(br domain.BundleRepository) *bundleService {
	return &bundleService{
		br: br,
	}
}

func (bs bundleService) FetchByID(ctx context.Context, bid int64) (*domain.Bundle, error) {
	return bs.br.GetByID(ctx, bid)
}

func (bs bundleService) FetchAll(ctx context.Context) (*[]domain.Bundle, error) {
	return bs.br.GetAll(ctx)
}

func (bs bundleService) Store(ctx context.Context, bundle *domain.Bundle, principal *domain.User) error {
	if !principal.HasClearance(domain.MEMBER) {
		return domain.NewNotAuthorizedErr("")
	}

	if bundle.ParentID < 0 {
		return domain.NewBadRequestErr("parent_id is invalid")
	}

	if bundle.ParentID > 0 {
		_, err := bs.br.GetByID(ctx, bundle.ParentID)
		if err != nil {
			return err
		}
	}

	return bs.br.Create(ctx, bundle)
}

func (bs bundleService) Remove(ctx context.Context, bid int64, principal *domain.User) error {
	if !principal.HasClearance(domain.MEMBER) {
		return domain.NewNotAuthorizedErr("")
	}

	if _, err := bs.br.GetByID(ctx, bid); err != nil {
		return err
	}

	return bs.br.Delete(ctx, bid)
}

func (bs bundleService) Update(ctx context.Context, bundle *domain.Bundle, principal *domain.User) error {
	if !principal.HasClearance(domain.MEMBER) {
		return domain.NewNotAuthorizedErr("")
	}

	if _, err := bs.br.GetByID(ctx, bundle.ID); err != nil {
		return err
	}

	return bs.br.Update(ctx, bundle)
}
