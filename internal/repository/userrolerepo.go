package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/go-sql-driver/mysql"
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

func (urr gormUserRoleRepository) GetByID(ctx context.Context, urid int64) (*domain.UserRole, error) {
	var role domain.UserRole
	res := urr.db.First(&role, urid)
	if err := res.Error; err != nil {
		switch {
		case errors.Is(gorm.ErrRecordNotFound, err):
			return nil, domain.NewRecordNotFoundErr("id", fmt.Sprint(urid))
		default:
			return nil, domain.NewInternalErr()
		}
	}

	return &role, nil
}

func (urr gormUserRoleRepository) GetAll(ctx context.Context) (*[]domain.UserRole, error) {
	var userroles []domain.UserRole
	res := urr.db.Find(&userroles)
	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &userroles, nil
}

func (urr gormUserRoleRepository) GetByUID(ctx context.Context, uid int64) (*[]domain.UserRole, error) {
	var userroles []domain.UserRole
	res := urr.db.Where("user_id = ?", uid).Find(&userroles)
	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &userroles, nil
}

func (urr gormUserRoleRepository) Create(ctx context.Context, userrole *domain.UserRole) error {
	res := urr.db.Create(userrole)
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
	}

	return nil
}

func (urr gormUserRoleRepository) CreateBatch(ctx context.Context, userroles *[]domain.UserRole) error {
	res := urr.db.Create(userroles)
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
	}

	return nil
}

func (urr gormUserRoleRepository) Update(ctx context.Context, userrole *domain.UserRole) error {
	res := urr.db.Updates(userrole)
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
	}

	return nil
}

func (urr gormUserRoleRepository) UpdateBatch(ctx context.Context, userroles *[]domain.UserRole) error {
	res := urr.db.Updates(userroles)
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
	}

	return nil
}

func (urr gormUserRoleRepository) Delete(ctx context.Context, rid int64) error {
	res := urr.db.Delete(&domain.UserRole{}, rid)
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (urr gormUserRoleRepository) DeleteBatch(ctx context.Context, rids []int64) error {
	res := urr.db.Delete(&domain.UserRole{}, rids)
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (urr gormUserRoleRepository) DeleteByRID(ctx context.Context, rid int64) error {
	res := urr.db.Where("role_id = ?", rid).Delete(&domain.UserRole{})
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (urr gormUserRoleRepository) DeleteByUID(ctx context.Context, uid int64) error {
	res := urr.db.Where("user_id = ?", uid).Delete(&domain.UserRole{})
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}
