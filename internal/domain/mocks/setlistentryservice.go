package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSetlistEntryService struct {
	mock.Mock
}

func (m MockSetlistEntryService) StoreBatch(ctx context.Context, entries *[]domain.SetlistEntry, principal *domain.User) error {
	ret := m.Called(ctx, entries, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSetlistEntryService) FetchByID(ctx context.Context, slid int64) (*domain.SetlistEntry, error) {
	ret := m.Called(ctx, slid)

	var r0 *domain.SetlistEntry
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.SetlistEntry)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockSetlistEntryService) FetchAll(ctx context.Context) (*[]domain.SetlistEntry, error) {
	ret := m.Called(ctx)

	var r0 *[]domain.SetlistEntry
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.SetlistEntry)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m MockSetlistEntryService) UpdateBatch(ctx context.Context, entries *[]domain.SetlistEntry, principal *domain.User) error {
	ret := m.Called(ctx, entries, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSetlistEntryService) RemoveBatch(ctx context.Context, setlist *domain.Setlist, ids []int64, principal *domain.User) error {
	ret := m.Called(ctx, setlist, ids, principal)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
