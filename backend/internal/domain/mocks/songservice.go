package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSongService struct {
	mock.Mock
}

func (m MockSongService) FetchByID(ctx context.Context, sid int64) (*domain.Song, error) {
	ret := m.Called(ctx, sid)

	var r0 *domain.Song
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Song)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
func (m MockSongService) FetchAll(ctx context.Context) (*[]domain.Song, error) {
	ret := m.Called(ctx)

	var r0 *[]domain.Song
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.Song)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockSongService) Fetch(ctx context.Context, options *domain.SongFilterOptions) ([]domain.Song, error) {
	ret := m.Called(ctx, options)

	var r0 []domain.Song
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]domain.Song)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockSongService) Update(ctx context.Context, song *domain.Song, principal *domain.User) error {
	ret := m.Called(ctx, song, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
func (m MockSongService) Store(ctx context.Context, song *domain.Song, principal *domain.User) error {
	ret := m.Called(ctx, song, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
func (m MockSongService) Remove(ctx context.Context, sid int64, principal *domain.User) error {
	ret := m.Called(ctx, sid, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
