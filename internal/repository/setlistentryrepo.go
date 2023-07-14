package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type gormSetlistEntryRepository struct {
	db *gorm.DB
}

//revive:disable:unexported-return
func NewGormSetlistEntryRepository(db *gorm.DB) *gormSetlistEntryRepository {
	return &gormSetlistEntryRepository{
		db: db,
	}
}

func (ser gormSetlistEntryRepository) GetByID(ctx context.Context, sid int64) (*domain.SetlistEntry, error) {
	var setlistEntry domain.SetlistEntry
	res := ser.db.First(&setlistEntry, sid)

	if err := res.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, domain.NewRecordNotFoundErr("id", fmt.Sprint(sid))
		default:
			return nil, domain.NewInternalErr()
		}
	}

	return &setlistEntry, nil
}

func (ser gormSetlistEntryRepository) GetAll(ctx context.Context) (*[]domain.SetlistEntry, error) {
	var setlists []domain.SetlistEntry
	res := ser.db.Find(&setlists)

	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &setlists, nil
}

func (ser gormSetlistEntryRepository) Create(ctx context.Context, setlistEntry *domain.SetlistEntry) error {
	res := ser.db.Create(setlistEntry)

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

func (ser gormSetlistEntryRepository) CreateBatch(ctx context.Context, setlistEntries *[]domain.SetlistEntry) error {
	res := ser.db.Create(setlistEntries)

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

func (ser gormSetlistEntryRepository) Delete(ctx context.Context, sid int64) error {
	setlistEntry := domain.SetlistEntry{ID: sid}
	res := ser.db.Delete(&setlistEntry)

	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (ser gormSetlistEntryRepository) DeleteBatch(ctx context.Context, sids []int64) error {
	setlistEntries := make([]domain.SetlistEntry, len(sids))
	for idx, val := range sids {
		setlistEntries[idx].ID = val
	}

	res := ser.db.Delete(&setlistEntries)

	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (ser gormSetlistEntryRepository) Update(ctx context.Context, setlistEntry *domain.SetlistEntry) error {
	res := ser.db.Updates(setlistEntry)

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

func (ser gormSetlistEntryRepository) UpdateBatch(ctx context.Context, setlistEntries *[]domain.SetlistEntry) error {
	res := ser.db.Updates(setlistEntries)

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
