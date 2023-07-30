package repository

import (
	"context"
	"errors"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type gormSetlistRoleRepository struct {
	db *gorm.DB
}

func (gsrs gormSetlistRoleRepository) Create(ctx context.Context, setlistRoles *[]domain.SetlistRole) error {
	res := gsrs.db.Create(setlistRoles)

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

func (gsrs gormSetlistRoleRepository) Get(ctx context.Context, setlistIDs []int64) (*[]domain.SetlistRole, error) {
	var retrievedSetlistRoles []domain.SetlistRole

	conditions := make(map[string][]int64, 0)

	if len(setlistIDs) <= 0 {
		conditions["setlist_id"] = setlistIDs
	}

	results := gsrs.db.Where(conditions).Find(&retrievedSetlistRoles)

	if err := results.Error; err != nil {
		return nil, nil
	}

	return &retrievedSetlistRoles, nil
}

func (gsrs gormSetlistRoleRepository) Update(ctx context.Context, setlistRoles *[]domain.SetlistRole) error {
	if setlistRoles == nil {
		return domain.NewInternalErr()
	}

	res := gsrs.db.Updates(setlistRoles)

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

func (gsrs gormSetlistRoleRepository) Delete(ctx context.Context, setlistRoleIDs []int64) error {
	setlistEntries := make([]domain.SetlistRole, len(setlistRoleIDs))
	for idx, val := range setlistRoleIDs {
		setlistEntries[idx].ID = val
	}

	res := gsrs.db.Delete(&setlistEntries)

	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}
