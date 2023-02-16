package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

type songService struct {
	sr domain.SongRepository
}

func NewSongService(sr domain.SongRepository) *songService {
	return &songService{
		sr: sr,
	}
}

func (ss songService) FetchByID(ctx context.Context, sid int64) (*domain.Song, error) {
	return nil, nil
}

func (ss songService) FetchAll(ctx context.Context) (*[]domain.Song, error) {
	return nil, nil
}

func (ss songService) Update(ctx context.Context, song *domain.Song, principal *domain.User) error {
	return nil
}

func (ss songService) Store(ctx context.Context, song *domain.Song, principal *domain.User) error {
	return nil
}

func (ss songService) Remove(ctx context.Context, sid int64, principal *domain.User) error {
	return nil
}
