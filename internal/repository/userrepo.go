package repository

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (ur userRepository) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (ur userRepository) CreateBatch(ctx context.Context, users []*domain.User) error {
	return nil
}

func (ur userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return nil, nil
}

func (ur userRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	return []*domain.User{}, nil
}
