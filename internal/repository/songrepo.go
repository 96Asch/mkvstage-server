package repository

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"gorm.io/gorm"
)

type gormSongRepository struct {
	db *gorm.DB
}

func NewGormSongRepository(db *gorm.DB) *gormSongRepository {
	return &gormSongRepository{
		db: db,
	}
}

func (sr gormSongRepository) GetByID(ctx context.Context, sid int64) (*domain.Song, error) {
	return nil, nil
}

func (sr gormSongRepository) GetAll(ctx context.Context) (*[]domain.Song, error) {
	return nil, nil
}

func (sr gormSongRepository) Create(ctx context.Context, song *domain.Song) error {
	return nil
}

func (sr gormSongRepository) Delete(ctx context.Context, sid int64) error {
	return nil
}

func (sr gormSongRepository) Update(ctx context.Context, song *domain.Song) error {
	return nil
}
