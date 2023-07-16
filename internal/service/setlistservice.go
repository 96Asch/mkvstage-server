package service

import (
	"context"
	"fmt"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

type setlistService struct {
	ur  domain.UserRepository
	slr domain.SetlistRepository
}

//revive:disable:unexported-return
func NewSetlistService(ur domain.UserRepository, slr domain.SetlistRepository) *setlistService {
	return &setlistService{
		ur:  ur,
		slr: slr,
	}
}

func (ss setlistService) FetchByID(ctx context.Context, sid int64) (*domain.Setlist, error) {
	setlist, err := ss.slr.GetByID(ctx, sid)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return setlist, nil
}

func (ss setlistService) FetchAll(ctx context.Context) (*[]domain.Setlist, error) {
	setlists, err := ss.slr.GetAll(ctx)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return setlists, nil
}

func (ss setlistService) Update(ctx context.Context, setlist *domain.Setlist, principal *domain.User) (*domain.Setlist, error) {
	currentSetlist, err := ss.slr.GetByID(ctx, setlist.ID)
	if err != nil {
		return nil, domain.FromError(err)
	}

	if !principal.HasClearance(domain.ADMIN) {
		if currentSetlist.CreatorID != principal.ID {
			return nil, domain.NewNotAuthorizedErr("Not authorized to update setlist")
		}
	}

	if setlist.Deadline.Before(time.Now()) {
		return nil, domain.NewBadRequestErr(fmt.Sprintf("%s must be later than %s", setlist.Deadline.String(), time.Now().String()))
	}

	if _, err := ss.ur.GetByID(ctx, setlist.CreatorID); err != nil {
		return nil, domain.FromError(err)
	}

	updatedSetlist, err := ss.slr.Update(ctx, setlist)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return updatedSetlist, nil
}

func (ss setlistService) Store(ctx context.Context, setlist *domain.Setlist, principal *domain.User) error {
	if !principal.HasClearance(domain.MEMBER) {
		return domain.NewNotAuthorizedErr("Not authorized to create setlists")
	}

	if setlist.Deadline.Before(time.Now()) {
		return domain.NewBadRequestErr(fmt.Sprintf("%s must be later than %s", setlist.Deadline.String(), time.Now().String()))
	}

	if _, err := ss.ur.GetByID(ctx, setlist.CreatorID); err != nil {
		return domain.FromError(err)
	}

	err := ss.slr.Create(ctx, setlist)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (ss setlistService) Remove(ctx context.Context, sid int64, principal *domain.User) error {
	if !principal.HasClearance(domain.ADMIN) {
		currentSetlist, err := ss.slr.GetByID(ctx, sid)
		if err != nil {
			return domain.FromError(err)
		}

		if currentSetlist.CreatorID != principal.ID {
			return domain.NewNotAuthorizedErr("user is neither an editor nor creator of the setlist")
		}
	}

	err := ss.slr.Delete(ctx, sid)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}
