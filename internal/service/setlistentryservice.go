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

		_, err := ses.sr.GetByID(ctx, entry.SongID)

		if err != nil {
			return domain.NewRecordNotFoundErr("SongID", fmt.Sprint(entry.SongID))
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
	return nil
}

func (ses setlistEntryService) RemoveBatch(ctx context.Context, id int64, principal *domain.User) error {
	return nil
}
