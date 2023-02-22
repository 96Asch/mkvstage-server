package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

type userRoleService struct {
	urr domain.UserRoleRepository
	rr  domain.RoleRepository
}

func NewUserRoleService(urr domain.UserRoleRepository, rr domain.RoleRepository) *userRoleService {
	return &userRoleService{
		urr: urr,
		rr:  rr,
	}
}

func containsUserRole(userRoles *[]domain.UserRole, urid int64) bool {
	for _, userRole := range *userRoles {
		if userRole.ID == urid {
			return true
		}
	}
	return false
}

func containsRole(roles *[]domain.Role, rid int64) bool {
	for _, role := range *roles {
		if role.ID == rid {
			return true
		}
	}
	return false
}

func (urs userRoleService) UpdateBatch(ctx context.Context, userRoles *[]domain.UserRole, principal *domain.User) error {
	currentUserRoles, err := urs.urr.GetByUID(ctx, principal.ID)
	if err != nil {
		return err
	}

	currentRoles, err := urs.rr.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, userrole := range *userRoles {
		if userrole.UserID != principal.ID {
			return domain.NewNotAuthorizedErr("cannot change other user's user roles")
		}

		if !containsUserRole(currentUserRoles, userrole.ID) {
			return domain.NewBadRequestErr("id is not valid")
		}

		if !containsRole(currentRoles, userrole.RoleID) {
			return domain.NewBadRequestErr("role_id is not valid")
		}
	}

	return urs.urr.UpdateBatch(ctx, userRoles)
}
