package repository

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"gorm.io/gorm"
)

type gormUserRoleRepository struct {
	db *gorm.DB
}

func NewGormUserRoleRepository(db *gorm.DB) *gormUserRoleRepository {
	return &gormUserRoleRepository{
		db: db,
	}
}

func (urr gormUserRoleRepository) GetByID(ctx context.Context, id int64) (*domain.UserRole, error) {
	return nil, nil
}

func (urr gormUserRoleRepository) GetAll(ctx context.Context) (*[]domain.UserRole, error) {
	return nil, nil
}

func (urr gormUserRoleRepository) Create(ctx context.Context, ur *domain.UserRole) error {
	return nil
}

func (urr gormUserRoleRepository) CreateBatch(ctx context.Context, urs *[]domain.UserRole) error {
	return nil
}

func (urr gormUserRoleRepository) Update(ctx context.Context, ur *domain.UserRole) error {
	return nil
}

func (urr gormUserRoleRepository) UpdateBatch(ctx context.Context, urs *[]domain.UserRole) error {
	return nil
}

func (urr gormUserRoleRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func (urr gormUserRoleRepository) DeleteBatch(ctx context.Context, ids []int64) error {
	return nil
}
