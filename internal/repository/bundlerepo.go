package repository

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"gorm.io/gorm"
)

type gormBundleRepository struct {
	DB *gorm.DB
}

func NewGormBundleRepository(db *gorm.DB) *gormBundleRepository {
	return &gormBundleRepository{
		DB: db,
	}
}

func (br gormBundleRepository) GetByID(ctx context.Context, bid int64) (*domain.Bundle, error) {
	return nil, nil
}
func (br gormBundleRepository) GetAll(ctx context.Context) (*[]domain.Bundle, error) {
	return nil, nil
}
func (br gormBundleRepository) Create(ctx context.Context, bundle *domain.Bundle) error {
	return nil
}
func (br gormBundleRepository) Delete(ctx context.Context, bid int64) error {
	return nil
}
func (br gormBundleRepository) Update(ctx context.Context, bundle *domain.Bundle) error {
	return nil
}
