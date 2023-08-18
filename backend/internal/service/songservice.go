package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/util"
)

type songService struct {
	ur domain.UserRepository
	sr domain.SongRepository
	br domain.BundleRepository
}

//revive:disable:unexported-return
func NewSongService(ur domain.UserRepository, sr domain.SongRepository, br domain.BundleRepository) *songService {
	return &songService{
		ur: ur,
		sr: sr,
		br: br,
	}
}

func (ss songService) FetchByID(ctx context.Context, sid int64) (*domain.Song, error) {
	song, err := ss.sr.GetByID(ctx, sid)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return song, nil
}

func (ss songService) FetchAll(ctx context.Context) (*[]domain.Song, error) {
	songs, err := ss.sr.GetAll(ctx)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return songs, nil
}

func (ss songService) Fetch(ctx context.Context, options *domain.SongFilterOptions) ([]domain.Song, error) {
	songs, err := ss.sr.Get(ctx, options)
	if err != nil {
		return nil, domain.FromError(err)
	}

	return songs, nil
}

func (ss songService) Update(ctx context.Context, song *domain.Song, principal *domain.User) error {
	if !principal.HasClearance(domain.EDITOR) {
		currentSong, err := ss.sr.GetByID(ctx, song.ID)
		if err != nil {
			return domain.FromError(err)
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

	if _, err := ss.br.GetByID(ctx, song.BundleID); err != nil {
		return domain.FromError(err)
	}

	if _, err := ss.ur.GetByID(ctx, song.CreatorID); err != nil {
		return domain.FromError(err)
	}

	err := ss.sr.Update(ctx, song)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (ss songService) Store(ctx context.Context, song *domain.Song, principal *domain.User) error {
	if !principal.HasClearance(domain.MEMBER) {
		return domain.NewNotAuthorizedErr("not authorized to create songs")
	}

	if song.CreatorID != principal.ID {
		return domain.NewBadRequestErr("cannot create a song with different creator")
	}

	if !song.IsValidKey() {
		return domain.NewBadRequestErr("invalid key")
	}

	if err := util.ValidateChordSheet(song.ChordSheet); err != nil {
		return domain.NewBadRequestErr(err.Error())
	}

	if _, err := ss.br.GetByID(ctx, song.BundleID); err != nil {
		return domain.FromError(err)
	}

	err := ss.sr.Create(ctx, song)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}

func (ss songService) Remove(ctx context.Context, sid int64, principal *domain.User) error {
	if !principal.HasClearance(domain.EDITOR) {
		currentSong, err := ss.sr.GetByID(ctx, sid)
		if err != nil {
			return domain.FromError(err)
		}

		if currentSong.CreatorID != principal.ID {
			return domain.NewNotAuthorizedErr("user is neither an editor nor creator of the song")
		}
	}

	err := ss.sr.Delete(ctx, sid)
	if err != nil {
		return domain.FromError(err)
	}

	return nil
}
