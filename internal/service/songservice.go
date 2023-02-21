package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
)

type songService struct {
	ur domain.UserRepository
	sr domain.SongRepository
}

func NewSongService(ur domain.UserRepository, sr domain.SongRepository) *songService {
	return &songService{
		ur: ur,
		sr: sr,
	}
}

func (ss songService) FetchByID(ctx context.Context, sid int64) (*domain.Song, error) {
	return ss.sr.GetByID(ctx, sid)
}

func (ss songService) FetchAll(ctx context.Context) (*[]domain.Song, error) {
	return ss.sr.GetAll(ctx)
}

func (ss songService) Update(ctx context.Context, song *domain.Song, principal *domain.User) error {
	if !principal.HasClearance(domain.EDITOR) {

		currentSong, err := ss.sr.GetByID(ctx, song.ID)
		if err != nil {
			return err
		}

		if currentSong.CreatorID != principal.ID {
			return domain.NewNotAuthorizedErr("user is neither an editor nor creator of the song")
		}

	}

	if !song.IsValidKey() {
		return domain.NewBadRequestErr("invalid key")
	}

	if err := util.ValidateChordSheet(song.ChordSheet); err != nil {
		return domain.NewBadRequestErr(err.Error())
	}

	_, err := ss.ur.GetByID(ctx, song.CreatorID)
	if err != nil {
		return err
	}

	return ss.sr.Update(ctx, song)
}

func (ss songService) Store(ctx context.Context, song *domain.Song, principal *domain.User) error {
	if !principal.HasClearance(domain.MEMBER) {
		return domain.NewNotAuthorizedErr("not authorized to create songs")
	}

	if !song.IsValidKey() {
		return domain.NewBadRequestErr("invalid key")
	}

	if err := util.ValidateChordSheet(song.ChordSheet); err != nil {
		return domain.NewBadRequestErr(err.Error())
	}

	if _, err := ss.ur.GetByID(ctx, song.CreatorID); err != nil {
		return err
	}

	return ss.sr.Create(ctx, song)
}

func (ss songService) Remove(ctx context.Context, sid int64, principal *domain.User) error {
	if !principal.HasClearance(domain.EDITOR) {

		currentSong, err := ss.sr.GetByID(ctx, sid)
		if err != nil {
			return err
		}

		if currentSong.CreatorID != principal.ID {
			return domain.NewNotAuthorizedErr("user is neither an editor nor creator of the song")
		}

	}

	return ss.sr.Delete(ctx, sid)
}
