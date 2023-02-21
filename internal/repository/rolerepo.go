package repository

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"gorm.io/gorm"
)

type gormRoleRepository struct {
	db *gorm.DB
}

func NewGormRoleRepository(db *gorm.DB) *gormRoleRepository {
	return &gormRoleRepository{
		db: db,
	}
}

func (gr gormRoleRepository) Create(ctx context.Context, role *domain.Role) error {
	return nil
}

func (gr gormRoleRepository) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	return nil, nil
}

func (gr gormRoleRepository) GetAll(ctx context.Context) (*[]domain.Role, error) {
	return nil, nil
}

func (gr gormRoleRepository) Update(ctx context.Context, role *domain.Role) error {
	return nil
}

func (gr gormRoleRepository) Deleter(ctx context.Context, rid int64) error {
	return nil
}
