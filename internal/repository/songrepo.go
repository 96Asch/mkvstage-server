package repository

import "gorm.io/gorm"

type gormSongRepository struct {
	db *gorm.DB
}

func NewGormSongRepository(db *gorm.DB) *gormSongRepository {
	return &gormSongRepository{
		db: db,
	}
}
