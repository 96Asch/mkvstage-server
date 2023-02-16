package service

import "github.com/96Asch/mkvstage-server/internal/domain"

type songService struct {
	sr domain.SongRepository
}

func NewSongService(sr domain.SongRepository) *songService {
	return &songService{
		sr: sr,
	}
}
