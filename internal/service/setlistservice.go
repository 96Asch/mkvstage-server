package service

import (
	"context"

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

func (ss setlistService) FetchAllGlobal(ctx context.Context, principal *domain.User) (*[]domain.Setlist, error) {
	setlists, err := ss.slr.GetAllGlobal(ctx, principal.ID)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return setlists, nil
}

func (ss setlistService) Update(ctx context.Context, setlist *domain.Setlist, principal *domain.User) (*domain.Setlist, error) {
	if !principal.HasClearance(domain.EDITOR) {
		currentSetlist, err := ss.slr.GetByID(ctx, setlist.ID)
		if err != nil {
			return nil, domain.FromError(err)
		}

		if currentSetlist.CreatorID != principal.ID {
			return nil, domain.NewNotAuthorizedErr("user is neither an editor nor creator of the setlist")
		}
	}

	_, err := ss.ur.GetByID(ctx, setlist.CreatorID)
	if err != nil {
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
		return domain.NewNotAuthorizedErr("not authorized to create setlists")
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
	if !principal.HasClearance(domain.EDITOR) {
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
