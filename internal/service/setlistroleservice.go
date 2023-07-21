package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

type setlistRoleService struct {
	slrr domain.SetlistRoleRepository
	slr  domain.SetlistRepository
	urr  domain.UserRoleRepository
}

func NewSetlistRoleService(slrr domain.SetlistRoleRepository, slr domain.SetlistRepository, urr domain.UserRoleRepository) *setlistRoleService {
	return &setlistRoleService{
		slrr: slrr,
		slr:  slr,
		urr:  urr,
	}
}

func (slrs setlistRoleService) Fetch(ctx context.Context, setlists *[]domain.Setlist) (*[]domain.SetlistRole, error) {
	setlistIDs := make([]int64, 0)

	if setlists != nil {
		for _, setlist := range *setlists {
			setlistIDs = append(setlistIDs, setlist.ID)
		}
	}

	retrievedSetlists, err := slrs.slrr.Get(ctx, setlistIDs)

	if err != nil {
		return nil, domain.FromError(err)
	}

	return retrievedSetlists, nil
}

func (slrs setlistRoleService) Store(ctx context.Context, setlistRoles *[]domain.SetlistRole, principal *domain.User) error {
	if principal == nil {
		return domain.NewNotAuthorizedErr("No user specified")
	}

	if setlistRoles == nil || len(*setlistRoles) <= 0 {
		return domain.NewBadRequestErr("No setlistroles given")
	}

	if !principal.HasClearance(domain.ADMIN) {
		userRoleIDs := make([]int64, len(*setlistRoles))

		for idx, setlistRole := range *setlistRoles {
			userRoleIDs[idx] = setlistRole.UserRoleID
		}

		retrievedUserRoles, userRoleErr := slrs.urr.Get(ctx, userRoleIDs)

		if userRoleErr != nil {
			return domain.FromError(userRoleErr)
		}

		for _, userrole := range *retrievedUserRoles {
			if principal.ID != userrole.UserID {
				return domain.NewNotAuthorizedErr("Cannot change the Setlist Role of someone else")
			}
		}
	}

	setlistIDs := make([]int64, len(*setlistRoles))

	for idx, setlistRole := range *setlistRoles {
		setlistIDs[idx] = setlistRole.SetlistID
	}

	if _, err := slrs.slr.GetByIDs(ctx, setlistIDs); err != nil {
		return domain.FromError(err)
	}

	err := slrs.slrr.Create(ctx, setlistRoles)

	if err != nil {
		return domain.FromError(err)
	}

	return nil

}

func (slrs setlistRoleService) Remove(ctx context.Context, setlistRoleIDs []int64, principal *domain.User) error {
	if principal == nil {
		return domain.NewNotAuthorizedErr("No user specified")
	}

	if len(setlistRoleIDs) <= 0 {
		return nil
	}

	retrievedSetlistRoles, setlistRoleErr := slrs.slrr.Get(ctx, setlistRoleIDs)

	if setlistRoleErr != nil {
		return domain.FromError(setlistRoleErr)
	}

	if !principal.HasClearance(domain.ADMIN) {
		userRoleIDs := make([]int64, len(*retrievedSetlistRoles))

		for idx, setlistRole := range *retrievedSetlistRoles {
			userRoleIDs[idx] = setlistRole.UserRoleID
		}

		retrievedUserRoles, userRoleErr := slrs.urr.Get(ctx, userRoleIDs)

		if userRoleErr != nil {
			return domain.FromError(userRoleErr)
		}

		for _, userrole := range *retrievedUserRoles {
			if userrole.UserID != principal.ID {
				return domain.NewNotAuthorizedErr("Cannot change the Setlist Role of someone else")
			}
		}
	}

	err := slrs.slrr.Delete(ctx, setlistRoleIDs)

	if err != nil {
		return domain.FromError(err)
	}

	return nil
}
