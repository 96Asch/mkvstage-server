package repository

import (
	"context"
	"errors"
	"log"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

//revive:disable:unexported-return
func NewGormUserRepository(db *gorm.DB) *gormUserRepository {
	return &gormUserRepository{
		db: db,
	}
}

func (ur gormUserRepository) Create(ctx context.Context, user *domain.User) error {
	res := ur.db.Create(user)
	if err := res.Error; err != nil {
		log.Println(err)
		return domain.NewInternalErr()
	}

	return nil
}

func (ur gormUserRepository) CreateBatch(ctx context.Context, users *[]domain.User) error {
	res := ur.db.CreateInBatches(users, 50)
	if res.Error != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (ur gormUserRepository) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	var user domain.User

	res := ur.db.First(&user, userID)
	if err := res.Error; err != nil {
		switch {
		case errors.Is(gorm.ErrRecordNotFound, err):
			return nil, domain.NewRecordNotFoundErr("id", string(rune(userID)))
		default:
			return nil, domain.NewInternalErr()
		}
	}

	return &user, nil
}

func (ur gormUserRepository) GetAll(ctx context.Context) (*[]domain.User, error) {
	var users []domain.User

	res := ur.db.Find(&users)
	if err := res.Error; err != nil {
		return nil, domain.NewInternalErr()
	}

	return &users, nil
}

func (ur gormUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	res := ur.db.Where("email = ?", &email).First(&user)
	if err := res.Error; err != nil {
		switch {
		case errors.Is(gorm.ErrRecordNotFound, err):
			return nil, domain.NewRecordNotFoundErr("email", email)
		default:
			return nil, domain.NewInternalErr()
		}
	}

	return &user, nil
}

// Update updates a user by the given non-zero user.ID and only updates columns
// with non-zero values.
func (ur gormUserRepository) Update(ctx context.Context, user *domain.User) error {
	res := ur.db.Updates(user)
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}

func (ur gormUserRepository) Delete(ctx context.Context, id int64) error {
	res := ur.db.Delete(&domain.User{}, id)
	if err := res.Error; err != nil {
		return domain.NewInternalErr()
	}

	return nil
}
