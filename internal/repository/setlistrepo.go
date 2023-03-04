package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type gormSetlistRepository struct {
	db *gorm.DB
}

//revive:disable:unexported-return
func NewGormSetlistRepository(db *gorm.DB) *gormSetlistRepository {
	return &gormSetlistRepository{
		db: db,
	}
}

func (sr gormSetlistRepository) GetByID(ctx context.Context, sid int64) (*domain.Setlist, error) {
	var setlist domain.Setlist
	res := sr.db.First(&setlist, sid)

	if err := res.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, domain.NewRecordNotFoundErr("id", fmt.Sprint(sid))
		default:
			return nil, domain.NewInternalErr()
		}
	}

	return &setlist, nil
}

func (sr gormSetlistRepository) GetAll(ctx context.Context) (*[]domain.Setlist, error) {
	var setlists []domain.Setlist
	res := sr.db.Find(&setlists)

	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &setlists, nil
}

func (sr gormSetlistRepository) GetAllGlobal(ctx context.Context, uid int64) (*[]domain.Setlist, error) {
	var setlists []domain.Setlist
	res := sr.db.Where("is_global = ?", true).Or("creator_id = ?", uid).Find(&setlists)

	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &setlists, nil
}

func (sr gormSetlistRepository) Create(ctx context.Context, setlist *domain.Setlist) error {
	res := sr.db.Create(setlist)

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

func (sr gormSetlistRepository) Delete(ctx context.Context, sid int64) error {
	setlist := domain.Setlist{ID: sid}
	res := sr.db.Delete(&setlist)

	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (sr gormSetlistRepository) Update(ctx context.Context, setlist *domain.Setlist) (*domain.Setlist, error) {
	res := sr.db.Updates(setlist)

	if err := res.Error; err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return nil, domain.NewBadRequestErr(mysqlErr.Message)
			default:
				return nil, domain.NewInternalErr()
			}
		}

		return nil, domain.NewInternalErr()
	}

	var updatedSetlist domain.Setlist

	res = sr.db.First(&updatedSetlist, setlist.ID)
	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &updatedSetlist, nil
}
