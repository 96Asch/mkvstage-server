package service

import (
	"context"
	"fmt"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
)

type setlistEntryService struct {
	sler domain.SetlistEntryRepository
	slr  domain.SetlistRepository
	sr   domain.SongRepository
}

//revive:disable:unexported-return
func NewSetlistEntryService(sler domain.SetlistEntryRepository, slr domain.SetlistRepository, sr domain.SongRepository) *setlistEntryService {
	return &setlistEntryService{
		sler: sler,
		slr:  slr,
		sr:   sr,
	}
}

func (ses setlistEntryService) StoreBatch(ctx context.Context, setlistEntries *[]domain.SetlistEntry, principal *domain.User) error {
	if principal == nil {
		return domain.NewInternalErr()
	}

	if !principal.HasClearance(domain.EDITOR) {
		return domain.NewNotAuthorizedErr("Invalid authorization")
	}

	if setlistEntries == nil {
		return domain.NewInternalErr()
	}

	if len(*setlistEntries) <= 0 {
		return nil
	}

	setlistID := (*setlistEntries)[0].SetlistID

	for _, entry := range *setlistEntries {
		if !util.IsValidTranpose(entry.Transpose) {
			return domain.NewBadRequestErr(fmt.Sprintf("Transpose must be between %d and %d", util.TransposeMin, util.TransposeMax))
		}

		if _, err := ses.sr.GetByID(ctx, entry.SongID); err != nil {
			return domain.FromError(err)
		}

		if setlistID != entry.SetlistID {
			return domain.NewBadRequestErr("SetlistID must be the same across entries")
		}
	}

	if _, err := ses.slr.GetByID(ctx, setlistID); err != nil {
		return domain.FromError(err)
	}

	err := ses.sler.CreateBatch(ctx, setlistEntries)

	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (ses setlistEntryService) FetchByID(ctx context.Context, id int64) (*domain.SetlistEntry, error) {
	setlistEntry, err := ses.sler.GetByID(ctx, id)

	if err != nil {
		return nil, domain.FromError(err)
	}

	return setlistEntry, nil
}

func (ses setlistEntryService) FetchAll(ctx context.Context) (*[]domain.SetlistEntry, error) {
	setlistEntries, err := ses.sler.GetAll(ctx)

	if err != nil {
		return nil, domain.FromError(err)
	}

	var minRank int64

	for _, entry := range *setlistEntries {
		if minRank > entry.Rank {
			return nil, domain.NewInternalErr()
		}

		minRank = entry.Rank
	}

	return setlistEntries, nil
}

func (ses setlistEntryService) FetchBySetlist(ctx context.Context, setlists *[]domain.Setlist) (*[]domain.SetlistEntry, error) {
	if setlists == nil {
		return nil, domain.NewInternalErr()
	}

	if len(*setlists) <= 0 {
		return nil, domain.NewBadRequestErr("No setlists given")
	}

	setlistEntries, err := ses.sler.GetBySetlist(ctx, setlists)

	if err != nil {
		return nil, domain.FromError(err)
	}

	var minRank int64

	for _, entry := range *setlistEntries {
		if minRank > entry.Rank {
			return nil, domain.NewInternalErr()
		}

		minRank = entry.Rank
	}

	return setlistEntries, nil
}

func (ses setlistEntryService) UpdateBatch(ctx context.Context, setlistEntries *[]domain.SetlistEntry, principal *domain.User) error {
	if principal == nil {
		return domain.NewInternalErr()
	}

	if !principal.HasClearance(domain.EDITOR) {
		return domain.NewNotAuthorizedErr("Invalid authorization")
	}

	if setlistEntries == nil {
		return domain.NewInternalErr()
	}

	setlistID := (*setlistEntries)[0].SetlistID

	for _, entry := range *setlistEntries {
		if !util.IsValidTranpose(entry.Transpose) {
			return domain.NewBadRequestErr(fmt.Sprintf("Transpose must be between %d and %d", util.TransposeMin, util.TransposeMax))
		}

		if _, err := ses.sr.GetByID(ctx, entry.SongID); err != nil {
			return domain.FromError(err)
		}

		if _, err := ses.sler.GetByID(ctx, entry.ID); err != nil {
			return domain.FromError(err)
		}

		if setlistID != entry.SetlistID {
			return domain.NewBadRequestErr("SetlistID must be the same across entries")
		}
	}

	if _, err := ses.slr.GetByID(ctx, setlistID); err != nil {
		return domain.FromError(err)
	}

	err := ses.sler.UpdateBatch(ctx, setlistEntries)

	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (ses setlistEntryService) RemoveBatch(ctx context.Context, setlist *domain.Setlist, ids []int64, principal *domain.User) error {
	if principal == nil {
		return domain.NewInternalErr()
	}

	if setlist == nil {
		return domain.NewInternalErr()
	}

	if !principal.HasClearance(domain.ADMIN) {
		if setlist.CreatorID != principal.ID {
			return domain.NewNotAuthorizedErr("Invalid authorization")
		}
	}

	if len(ids) <= 0 {
		return nil
	}

	for _, id := range ids {
		if _, err := ses.sler.GetByID(ctx, id); err != nil {
			return domain.FromError(err)
		}
	}

	err := ses.sler.DeleteBatch(ctx, ids)

	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (ses setlistEntryService) RemoveBySetlist(ctx context.Context, setlist *domain.Setlist, principal *domain.User) error {
	if principal == nil {
		return domain.NewInternalErr()
	}

	if setlist == nil {
		return domain.NewInternalErr()
	}

	if !principal.HasClearance(domain.ADMIN) {
		if setlist.CreatorID != principal.ID {
			return domain.NewNotAuthorizedErr("Invalid authorization")
		}
	}

	setlists := &[]domain.Setlist{*setlist}
	toDeleteSetlistEntries, err := ses.sler.GetBySetlist(ctx, setlists)

	if err != nil {
		return domain.FromError(err)
	}

	toDeleteSetlistEntryIDs := make([]int64, len(*toDeleteSetlistEntries))

	for idx, entry := range *toDeleteSetlistEntries {
		toDeleteSetlistEntryIDs[idx] = entry.ID
	}

	if err := ses.sler.DeleteBatch(ctx, toDeleteSetlistEntryIDs); err != nil {
		return domain.FromError(err)
	}

	return nil
}
