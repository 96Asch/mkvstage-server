package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"gorm.io/gorm"
)

type gormBundleRepository struct {
	db *gorm.DB
}

func NewGormBundleRepository(db *gorm.DB) *gormBundleRepository {
	return &gormBundleRepository{
		db: db,
	}
}

func (br gormBundleRepository) GetByID(ctx context.Context, bid int64) (*domain.Bundle, error) {
	var bundle domain.Bundle
	res := br.db.First(&bundle, bid)
	if err := res.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, domain.NewRecordNotFoundErr("id", fmt.Sprint(bid))
		default:
			return nil, domain.NewInternalErr()
		}
	}

	return &bundle, nil
}
func (br gormBundleRepository) GetAll(ctx context.Context) (*[]domain.Bundle, error) {
	var bundles []domain.Bundle
	res := br.db.Find(&bundles)
	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &bundles, nil
}

func (br gormBundleRepository) GetLeaves(ctx context.Context) (*[]domain.Bundle, error) {

	var bundles []domain.Bundle

	res := br.db.Unscoped().
		Table("bundles b").
		Where("NOT EXISTS (?) AND deleted_at IS NULL",
			br.db.Unscoped().
				Model(&domain.Bundle{}).
				Select("NULL").
				Where("parent_id = b.id")).
		Find(&bundles)

	if err := res.Error; err != nil {
		log.Println(err)
		return nil, domain.NewInternalErr()
	}

	return &bundles, nil
}

func (br gormBundleRepository) Create(ctx context.Context, bundle *domain.Bundle) error {
	res := br.db.Create(bundle)
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (br gormBundleRepository) Delete(ctx context.Context, bid int64) error {
	bundle := domain.Bundle{ID: bid}
	res := br.db.Delete(&bundle)
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}
func (br gormBundleRepository) Update(ctx context.Context, bundle *domain.Bundle) error {
	res := br.db.Save(bundle)
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}
