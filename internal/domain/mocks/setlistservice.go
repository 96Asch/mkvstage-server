package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSetlistService struct {
	mock.Mock
}

func (m MockSetlistService) FetchByID(ctx context.Context, slid int64) (*domain.Setlist, error) {
	ret := m.Called(ctx, slid)

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
func (m MockSetlistService) FetchAll(ctx context.Context) (*[]domain.Setlist, error) {
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

func (m MockSetlistService) FetchAllGlobal(ctx context.Context, principal *domain.User) (*[]domain.Setlist, error) {
	ret := m.Called(ctx, principal)

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

func (m MockSetlistService) Update(ctx context.Context, setlist *domain.Setlist, principal *domain.User) (*domain.Setlist, error) {
	ret := m.Called(ctx, setlist, principal)

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

func (m MockSetlistService) Store(ctx context.Context, setlist *domain.Setlist, principal *domain.User) error {
	ret := m.Called(ctx, setlist, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSetlistService) Remove(ctx context.Context, slid int64, principal *domain.User) error {
	ret := m.Called(ctx, slid, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
