package repository

import (
	"context"
	"errors"

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
	res := ur.db.Create(user)
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}
	return nil
}

func (ur userRepository) CreateBatch(ctx context.Context, users []*domain.User) error {
	res := ur.db.CreateInBatches(users, 50)
	if res.Error != nil {
		return domain.NewInternalErr()
	}
	return nil
}

func (ur userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	res := ur.db.First(&user, id)
	if err := res.Error; err != nil {

		switch {
		case errors.Is(gorm.ErrRecordNotFound, err):
			return nil, domain.NewRecordNotFoundErr("id", string(rune(id)))
		default:
			return nil, domain.NewInternalErr()
		}

	}

	return nil, nil
}

func (ur userRepository) GetAll(ctx context.Context) (*[]domain.User, error) {
	var users []domain.User
	res := ur.db.Find(&users)
	if err := res.Error; err != nil {
		return nil, nil
	}

	return &users, nil
}
