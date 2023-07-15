package service

import (
	"context"
	"fmt"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
)

type setlistEntryService struct {
	ser domain.SetlistEntryRepository
	sr  domain.SongRepository
}

//revive:disable:unexported-return
func NewSetlistEntryService(ser domain.SetlistEntryRepository, sr domain.SongRepository) *setlistEntryService {
	return &setlistEntryService{
		ser: ser,
		sr:  sr,
	}
}

func (ses setlistEntryService) StoreBatch(ctx context.Context, setlistEntries *[]domain.SetlistEntry, principal *domain.User) error {
	if !principal.HasClearance(domain.EDITOR) {
		return domain.NewNotAuthorizedErr("Invalid authorization")
	}

	for _, entry := range *setlistEntries {
		if !util.IsValidTranpose(entry.Transpose) {
			return domain.NewBadRequestErr(fmt.Sprintf("Transpose must be between %d and %d", util.TransposeMin, util.TransposeMax))
		}

		if _, err := ses.sr.GetByID(ctx, entry.SongID); err != nil {
			return domain.FromError(err)
		}
	}

	err := ses.ser.CreateBatch(ctx, setlistEntries)

	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (ses setlistEntryService) FetchByID(ctx context.Context, id int64) (*domain.SetlistEntry, error) {
	setlistEntry, err := ses.ser.GetByID(ctx, id)

	if err != nil {
		return nil, domain.FromError(err)
	}

	return setlistEntry, nil
}

func (ses setlistEntryService) FetchAll(ctx context.Context) (*[]domain.SetlistEntry, error) {
	setlistEntries, err := ses.ser.GetAll(ctx)

	if err != nil {
		return nil, domain.FromError(err)
	}

	return setlistEntries, nil
}

func (ses setlistEntryService) UpdateBatch(ctx context.Context, setlistEntries *[]domain.SetlistEntry, principal *domain.User) error {
	if !principal.HasClearance(domain.EDITOR) {
		return domain.NewNotAuthorizedErr("Invalid authorization")
	}

	for _, entry := range *setlistEntries {
		if !util.IsValidTranpose(entry.Transpose) {
			return domain.NewBadRequestErr(fmt.Sprintf("Transpose must be between %d and %d", util.TransposeMin, util.TransposeMax))
		}

		if _, err := ses.sr.GetByID(ctx, entry.SongID); err != nil {
			return domain.FromError(err)
		}

		if _, err := ses.ser.GetByID(ctx, entry.ID); err != nil {
			return domain.FromError(err)
		}
	}

	err := ses.ser.UpdateBatch(ctx, setlistEntries)

	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (ses setlistEntryService) RemoveBatch(ctx context.Context, ids []int64, principal *domain.User) error {
	if !principal.HasClearance(domain.EDITOR) {
		return domain.NewNotAuthorizedErr("Invalid authorization")
	}

	for _, id := range ids {
		if _, err := ses.ser.GetByID(ctx, id); err != nil {
			return domain.FromError(err)
		}
	}

	err := ses.ser.DeleteBatch(ctx, ids)

	if err != nil {
		return domain.FromError(err)
	}

	return nil
}
