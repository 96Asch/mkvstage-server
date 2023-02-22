package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
)

type userService struct {
	ur  domain.UserRepository
	rr  domain.RoleRepository
	urr domain.UserRoleRepository
}

func NewUserService(ur domain.UserRepository, rr domain.RoleRepository, urr domain.UserRoleRepository) domain.UserService {
	return &userService{
		ur:  ur,
		rr:  rr,
		urr: urr,
	}
}

func (us *userService) FetchByID(ctx context.Context, id int64) (*domain.User, error) {
	return us.ur.GetByID(ctx, id)
}

func (us *userService) FetchAll(ctx context.Context) (*[]domain.User, error) {

	users, err := us.ur.GetAll(ctx)
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

	err := us.ur.Create(ctx, user)
	if err != nil {
		return err
	}

	currentRoles, err := us.rr.GetAll(ctx)
	if err != nil {
		return err
	}

	userroles := make([]domain.UserRole, len(*currentRoles))
	for idx, role := range *currentRoles {
		userroles[idx] = domain.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}
	}

	err = us.urr.CreateBatch(ctx, &userroles)
	if err != nil {
		return err
	}

	return nil
}

func (us userService) Update(ctx context.Context, user *domain.User) error {

	if user.ID == 0 {
		return domain.NewRecordNotFoundErr("user_id", "0")
	}

	return us.ur.Update(ctx, user)
}

func (us userService) Remove(ctx context.Context, user *domain.User, id int64) (int64, error) {

	deleteId := id
	if id == 0 {
		deleteId = user.ID
	}

	if user.ID != deleteId && !user.HasClearance(domain.ADMIN) {
		return 0, domain.NewNotAuthorizedErr("cannot delete given id")
	}

	if _, err := us.ur.GetByID(ctx, deleteId); err != nil {
		return 0, err
	}

	if err := us.ur.Delete(ctx, deleteId); err != nil {
		return 0, err
	}

	if err := us.urr.DeleteByUID(ctx, deleteId); err != nil {
		return 0, err
	}

	return deleteId, nil
}

func (us userService) Authorize(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := us.ur.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := util.Validate(password, user.Password); err != nil {
		return nil, domain.NewNotAuthorizedErr("email and/or password is incorrect")
	}

	return user, nil
}
