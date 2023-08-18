package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSongRepository struct {
	mock.Mock
}

func (m MockSongRepository) GetByID(ctx context.Context, sid int64) (*domain.Song, error) {
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

func (m MockSongRepository) GetAll(ctx context.Context) (*[]domain.Song, error) {
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

func (m MockSongRepository) Create(ctx context.Context, song *domain.Song) error {
	ret := m.Called(ctx, song)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSongRepository) Delete(ctx context.Context, sid int64) error {
	ret := m.Called(ctx, sid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSongRepository) Update(ctx context.Context, song *domain.Song) error {
	ret := m.Called(ctx, song)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
