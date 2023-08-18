package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type gormSongRepository struct {
	db *gorm.DB
}

//revive:disable:unexported-return
func NewGormSongRepository(db *gorm.DB) *gormSongRepository {
	return &gormSongRepository{
		db: db,
	}
}

func (sr gormSongRepository) GetByID(ctx context.Context, sid int64) (*domain.Song, error) {
	var song domain.Song
	res := sr.db.First(&song, sid)

	if err := res.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, domain.NewRecordNotFoundErr("id", fmt.Sprint(sid))
		default:
			return nil, domain.NewInternalErr()
		}
	}

	return &song, nil
}

func (sr gormSongRepository) GetAll(ctx context.Context) (*[]domain.Song, error) {
	var songs []domain.Song
	res := sr.db.Find(&songs)

	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &songs, nil
}

func (sr gormSongRepository) Get(ctx context.Context, options *domain.SongFilterOptions) ([]domain.Song, error) {
	return []domain.Song{}, nil
}

func (sr gormSongRepository) Create(ctx context.Context, song *domain.Song) error {
	res := sr.db.Create(song)

	if err := res.Error; err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return domain.NewBadRequestErr(mysqlErr.Message)
			default:
				return domain.NewInternalErr()
			}
		}

		return domain.NewInternalErr()
	}

	return nil
}

func (sr gormSongRepository) Delete(ctx context.Context, sid int64) error {
	song := domain.Song{ID: sid}
	res := sr.db.Delete(&song)

	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (sr gormSongRepository) Update(ctx context.Context, song *domain.Song) error {
	res := sr.db.Updates(song)

	if err := res.Error; err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return domain.NewBadRequestErr(mysqlErr.Message)
			default:
				return domain.NewInternalErr()
			}
		}

		return domain.NewInternalErr()
	}

	return nil
}
