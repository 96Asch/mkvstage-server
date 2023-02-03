package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(ur domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: ur,
	}
}

func (us *userService) FetchByID(ctx context.Context, id int64) (*domain.User, error) {
	return us.userRepo.GetByID(ctx, id)
}

func (us *userService) FetchAll(ctx context.Context) (*[]domain.PublicUser, error) {

	users, err := us.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	publicUsers := make([]domain.PublicUser, len(*users))
	for idx, user := range *users {
		publicUsers[idx] = domain.PublicUser{
			ID:           user.ID,
			LastName:     user.LastName,
			FirstName:    user.FirstName,
			Email:        user.Email,
			Permission:   user.Permission,
			ProfileColor: user.ProfileColor,
			UpdatedAt:    user.UpdatedAt,
		}

	}

	return &publicUsers, nil
}

func (us userService) Store(ctx context.Context, user *domain.User) error {
	return us.userRepo.Create(ctx, user)
}
