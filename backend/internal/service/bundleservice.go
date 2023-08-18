package service

import (
	"context"
	"log"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
)

type bundleService struct {
	br domain.BundleRepository
}

//revive:disable:unexported-return
func NewBundleService(br domain.BundleRepository) *bundleService {
	return &bundleService{
		br: br,
	}
}

func (bs bundleService) FetchByID(ctx context.Context, bid int64) (*domain.Bundle, error) {
	bundle, err := bs.br.GetByID(ctx, bid)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return bundle, nil
}

func (bs bundleService) FetchAll(ctx context.Context) (*[]domain.Bundle, error) {
	bundles, err := bs.br.GetAll(ctx)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return bundles, nil
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
			return domain.FromError(err)
		}
	}

	err := bs.br.Create(ctx, bundle)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func contains(bundles *[]domain.Bundle, bid int64) bool {
	for _, bundle := range *bundles {
		if bundle.ID == bid {
			return true
		}
	}

	return false
}

func (bs bundleService) Remove(ctx context.Context, bid int64, principal *domain.User) error {
	if !principal.HasClearance(domain.MEMBER) {
		return domain.NewNotAuthorizedErr("")
	}

	if _, err := bs.br.GetByID(ctx, bid); err != nil {
		return domain.FromError(err)
	}

	leaves, err := bs.br.GetLeaves(ctx)
	if err != nil {
		return domain.FromError(err)
	}

	log.Println(leaves)

	if !contains(leaves, bid) {
		return domain.NewBadRequestErr("given id is not a leaf bundle")
	}

	err = bs.br.Delete(ctx, bid)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (bs bundleService) Update(ctx context.Context, bundle *domain.Bundle, principal *domain.User) error {
	if !principal.HasClearance(domain.MEMBER) {
		return domain.NewNotAuthorizedErr("")
	}

	if _, err := bs.br.GetByID(ctx, bundle.ID); err != nil {
		return domain.FromError(err)
	}

	err := bs.br.Update(ctx, bundle)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}
