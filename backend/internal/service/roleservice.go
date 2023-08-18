package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
)

type roleService struct {
	rr  domain.RoleRepository
	ur  domain.UserRepository
	urr domain.UserRoleRepository
}

//revive:disable:unexported-return
func NewRoleService(rr domain.RoleRepository, ur domain.UserRepository, urr domain.UserRoleRepository) *roleService {
	return &roleService{
		rr:  rr,
		ur:  ur,
		urr: urr,
	}
}

func (rs roleService) FetchByID(ctx context.Context, rid int64) (*domain.Role, error) {
	role, err := rs.rr.GetByID(ctx, rid)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return role, nil
}

func (rs roleService) FetchAll(ctx context.Context) (*[]domain.Role, error) {
	role, err := rs.rr.GetAll(ctx)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return role, nil
}

func (rs roleService) Update(ctx context.Context, role *domain.Role, principal *domain.User) error {
	if role.ID == 0 {
		return domain.NewBadRequestErr("id cannot be zero")
	}

	if principal.Permission != domain.ADMIN {
		return domain.NewNotAuthorizedErr("not authorized to update roles")
	}

	_, err := rs.rr.GetByID(ctx, role.ID)
	if err != nil {
		return domain.FromError(err)
	}

	err = rs.rr.Update(ctx, role)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (rs roleService) Store(ctx context.Context, role *domain.Role, principal *domain.User) error {
	if principal.Permission != domain.ADMIN {
		return domain.NewNotAuthorizedErr("not authorized to create roles")
	}

	err := rs.rr.Create(ctx, role)
	if err != nil {
		return domain.FromError(err)
	}

	users, err := rs.ur.GetAll(ctx)
	if err != nil {
		return domain.FromError(err)
	}

	userroles := make([]domain.UserRole, len(*users))
	for idx, user := range *users {
		userroles[idx] = domain.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}
	}

	err = rs.urr.CreateBatch(ctx, &userroles)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (rs roleService) Remove(ctx context.Context, rid int64, principal *domain.User) error {
	if principal.Permission != domain.ADMIN {
		return domain.NewNotAuthorizedErr("not authorized to create roles")
	}

	err := rs.rr.Delete(ctx, rid)
	if err != nil {
		return domain.FromError(err)
	}

	err = rs.urr.DeleteByRID(ctx, rid)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}
