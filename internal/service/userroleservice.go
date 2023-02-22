package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

type userRoleService struct {
	urr domain.UserRoleRepository
}

func NewUserRoleService(urr domain.UserRoleRepository) *userRoleService {
	return &userRoleService{
		urr: urr,
	}
}

func (urs userRoleService) GetAll(ctx context.Context) (*[]domain.UserRole, error) {
	return urs.urr.GetAll(ctx)
}

func (urs userRoleService) GetByUser(ctx context.Context, user *domain.User) (*[]domain.UserRole, error) {
	return urs.urr.GetByUID(ctx, user.ID)
}

func containsUserRole(userRoles *[]domain.UserRole, urid int64) (int, bool) {
	for idx, userRole := range *userRoles {
		if userRole.ID == urid {
			return idx, true
		}
	}
	return 0, false
}

func (urs userRoleService) SetActiveBatch(ctx context.Context, urids []int64, principal *domain.User) (*[]domain.UserRole, error) {
	userRoles, err := urs.urr.GetByUID(ctx, principal.ID)
	if err != nil {
		return nil, err
	}

	toUpdateUserRoles := make([]domain.UserRole, len(urids))
	for idx, urid := range urids {
		currentUserRoleIdx, exists := containsUserRole(userRoles, urid)
		if !exists {
			return nil, domain.NewBadRequestErr("invalid id given")
		}

		(*userRoles)[currentUserRoleIdx].Active = true
		toUpdateUserRoles[idx] = (*userRoles)[currentUserRoleIdx]
	}

	err = urs.urr.UpdateBatch(ctx, &toUpdateUserRoles)
	if err != nil {
		return nil, err
	}

	return &toUpdateUserRoles, nil
}
