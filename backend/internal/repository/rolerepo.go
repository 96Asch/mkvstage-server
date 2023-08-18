package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type gormRoleRepository struct {
	db *gorm.DB
}

//revive:disable:unexported-return
func NewGormRoleRepository(db *gorm.DB) *gormRoleRepository {
	return &gormRoleRepository{
		db: db,
	}
}

func (rr gormRoleRepository) Create(ctx context.Context, role *domain.Role) error {
	res := rr.db.Create(role)
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

func (rr gormRoleRepository) GetByID(ctx context.Context, rid int64) (*domain.Role, error) {
	var role domain.Role

	res := rr.db.First(&role, rid)
	if err := res.Error; err != nil {
		switch {
		case errors.Is(gorm.ErrRecordNotFound, err):
			return nil, domain.NewRecordNotFoundErr("id", fmt.Sprint(rid))
		default:
			return nil, domain.NewInternalErr()
		}
	}

	return &role, nil
}

func (rr gormRoleRepository) GetAll(ctx context.Context) (*[]domain.Role, error) {
	var roles []domain.Role

	res := rr.db.Find(&roles)
	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &roles, nil
}

func (rr gormRoleRepository) Update(ctx context.Context, role *domain.Role) error {
	res := rr.db.Updates(role)
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

func (rr gormRoleRepository) Delete(ctx context.Context, rid int64) error {
	res := rr.db.Delete(&domain.Role{}, rid)
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}
