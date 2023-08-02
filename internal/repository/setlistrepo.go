package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (slr gormSetlistRepository) GetByID(ctx context.Context, sid int64) (*domain.Setlist, error) {
	var setlist domain.Setlist
	res := slr.db.First(&setlist, sid)

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

func (slr gormSetlistRepository) GetByIDs(ctx context.Context, sids []int64) (*[]domain.Setlist, error) {
	var setlists []domain.Setlist
	res := slr.db.First(&setlists, sids)

	if err := res.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, domain.NewRecordNotFoundErr("id", fmt.Sprint(sids))
		default:
			return nil, domain.NewInternalErr()
		}
	}

	return &setlists, nil
}

func (slr gormSetlistRepository) GetAll(ctx context.Context) (*[]domain.Setlist, error) {
	var setlists []domain.Setlist
	res := slr.db.Find(&setlists)

	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &setlists, nil
}

func (slr gormSetlistRepository) Get(ctx context.Context, fromTime time.Time, toTime time.Time) (*[]domain.Setlist, error) {
	var setlists []domain.Setlist

	var result *gorm.DB

	if fromTime.IsZero() || toTime.IsZero() {
		result = slr.db.
			Where("deadline BETWEEN ? AND ?", fromTime, toTime).
			Find(&setlists)
	} else {
		result = slr.db.
			Find(&setlists)
	}

	if err := result.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &setlists, nil
}

func (slr gormSetlistRepository) GetByTimeframe(ctx context.Context, from time.Time, to time.Time) (*[]domain.Setlist, error) {
	var setlists []domain.Setlist
	res := slr.db.
		Where("deadline BETWEEN ? AND ?", from, to).
		Find(&setlists)

	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &setlists, nil
}

func (slr gormSetlistRepository) Create(ctx context.Context, setlist *domain.Setlist) error {
	res := slr.db.Create(setlist)

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

func (slr gormSetlistRepository) Delete(ctx context.Context, sid int64) error {
	setlist := domain.Setlist{ID: sid}
	res := slr.db.Delete(&setlist)

	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (slr gormSetlistRepository) Update(ctx context.Context, setlist *domain.Setlist) (*domain.Setlist, error) {
	res := slr.db.Updates(setlist)

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

	res = slr.db.First(&updatedSetlist, setlist.ID)
	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &updatedSetlist, nil
}
