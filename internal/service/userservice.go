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
	user, err := us.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) FetchAll(ctx context.Context) ([]*domain.User, error) {
	return []*domain.User{}, nil
}

func (us userService) Store(ctx context.Context, user *domain.User) error {
	return nil
}
