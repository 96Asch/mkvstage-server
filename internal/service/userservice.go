package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
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

	if user.Password != "" {
		password := user.Password

		hash, err := util.Encrypt(password)
		if err != nil {
			return domain.NewInternalErr()
		}

		user.Password = hash
	}

	return us.userRepo.Create(ctx, user)
}

func (us userService) Update(ctx context.Context, user *domain.User) error {

	if user.ID == 0 {
		return domain.NewRecordNotFoundErr("user_id", "0")
	}

	return us.userRepo.Update(ctx, user)
}

func (us userService) Remove(ctx context.Context, user *domain.User, id int64) (int64, error) {

	deleteId := id
	if id == 0 {
		deleteId = user.ID
	}

	if user.ID != deleteId && !user.HasClearance() {
		return 0, domain.NewNotAuthorizedErr("cannot delete given id")
	}

	if _, err := us.userRepo.GetByID(ctx, deleteId); err != nil {
		return 0, err
	}

	if err := us.userRepo.Delete(ctx, deleteId); err != nil {
		return 0, err
	}

	return deleteId, nil
}

func (us userService) Authorize(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := us.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := util.Validate(password, user.Password); err != nil {
		return nil, domain.NewNotAuthorizedErr("email and/or password is incorrect")
	}

	return user, nil
}
