package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSetlistRepository struct {
	mock.Mock
}

func (m MockSetlistRepository) GetByID(ctx context.Context, sid int64) (*domain.Setlist, error) {
	ret := m.Called(ctx, sid)

	var r0 *domain.Setlist
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Setlist)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockSetlistRepository) GetAll(ctx context.Context) (*[]domain.Setlist, error) {
	ret := m.Called(ctx)

	var r0 *[]domain.Setlist
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.Setlist)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockSetlistRepository) Create(ctx context.Context, setlist *domain.Setlist) error {
	ret := m.Called(ctx, setlist)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSetlistRepository) Delete(ctx context.Context, sid int64) error {
	ret := m.Called(ctx, sid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSetlistRepository) Update(ctx context.Context, setlist *domain.Setlist) (*domain.Setlist, error) {
	ret := m.Called(ctx, setlist)

	var r0 *domain.Setlist
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Setlist)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
