package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
)

type userRoleService struct {
	urr domain.UserRoleRepository
}

//revive:disable:unexported-return
func NewUserRoleService(urr domain.UserRoleRepository) *userRoleService {
	return &userRoleService{
		urr: urr,
	}
}

func (urs userRoleService) FetchAll(ctx context.Context) (*[]domain.UserRole, error) {
	userrole, err := urs.urr.GetAll(ctx)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return userrole, nil
}

func (urs userRoleService) FetchByUser(ctx context.Context, user *domain.User) (*[]domain.UserRole, error) {
	userroles, err := urs.urr.GetByUID(ctx, user.ID)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return userroles, nil
}

func containsUserRole(userRoles *[]domain.UserRole, urid int64) bool {
	for _, userRole := range *userRoles {
		if userRole.ID == urid {
			return true
		}
	}

	return false
}

func containsID(ids []int64, id int64) bool {
	for _, _id := range ids {
		if _id == id {
			return true
		}
	}

	return false
}

func (urs userRoleService) SetActiveBatch(ctx context.Context, urids []int64, principal *domain.User) (*[]domain.UserRole, error) {
	userRoles, err := urs.urr.GetByUID(ctx, principal.ID)
	if err != nil {
		return nil, domain.FromError(err)
	}

	for _, urid := range urids {
		if !containsUserRole(userRoles, urid) {
			return nil, domain.NewBadRequestErr("invalid id given")
		}
	}

	toUpdateUserRoles := make([]domain.UserRole, 0)

	for _, userrole := range *userRoles {
		if !userrole.Active && containsID(urids, userrole.ID) {
			toUpdateUserRoles = append(toUpdateUserRoles, domain.UserRole{
				ID:     userrole.ID,
				UserID: userrole.UserID,
				User:   userrole.User,
				RoleID: userrole.RoleID,
				Role:   userrole.Role,
				Active: true,
			})
		} else if userrole.Active && !containsID(urids, userrole.ID) {
			toUpdateUserRoles = append(toUpdateUserRoles, domain.UserRole{
				ID:     userrole.ID,
				UserID: userrole.UserID,
				User:   userrole.User,
				RoleID: userrole.RoleID,
				Role:   userrole.Role,
				Active: false,
			})
		}
	}

	if len(toUpdateUserRoles) == 0 {
		return nil, domain.NewBadRequestErr("no changes were made")
	}

	err = urs.urr.UpdateBatch(ctx, &toUpdateUserRoles)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return &toUpdateUserRoles, nil
}
