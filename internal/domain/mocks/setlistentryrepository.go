package mocks

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSetlistEntryRepository struct {
	mock.Mock
}

func (m MockSetlistEntryRepository) GetByID(ctx context.Context, sid int64) (*domain.SetlistEntry, error) {
	ret := m.Called(ctx, sid)

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

func (m MockSetlistEntryRepository) GetAll(ctx context.Context) (*[]domain.SetlistEntry, error) {
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

func (m MockSetlistEntryRepository) Create(ctx context.Context, setlistEntry *domain.SetlistEntry) error {
	ret := m.Called(ctx, setlistEntry)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSetlistEntryRepository) CreateBatch(ctx context.Context, setlistEntries *[]domain.SetlistEntry) error {
	ret := m.Called(ctx, setlistEntries)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSetlistEntryRepository) Delete(ctx context.Context, sid int64) error {
	ret := m.Called(ctx, sid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSetlistEntryRepository) DeleteBatch(ctx context.Context, sids []int64) error {
	ret := m.Called(ctx, sids)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSetlistEntryRepository) Update(ctx context.Context, setlistEntry *domain.SetlistEntry) error {
	ret := m.Called(ctx, setlistEntry)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m MockSetlistEntryRepository) UpdateBatch(ctx context.Context, setlistEntries *[]domain.SetlistEntry) error {
	ret := m.Called(ctx, setlistEntries)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
