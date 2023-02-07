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

func (us *userService) FetchAll(ctx context.Context) (*[]domain.User, error) {

	users, err := us.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us userService) Store(ctx context.Context, user *domain.User) error {
	return us.userRepo.Create(ctx, user)
}

func (us userService) Update(ctx context.Context, user *domain.User) error {

	if user.ID == 0 {
		return domain.NewRecordNotFoundErr("user_id", "0")
	}

	return us.userRepo.Update(ctx, user)
}

func (us userService) Remove(ctx context.Context, user *domain.User, id int64) error {

	deleteId := id
	if id == 0 {
		deleteId = user.ID
	}

	if user.ID != id && !user.HasClearance() {
		return domain.NewNotAuthorizedErr("cannot delete given id")
	}

	if err := us.userRepo.Delete(ctx, deleteId); err != nil {
		return err
	}

	return nil
}
