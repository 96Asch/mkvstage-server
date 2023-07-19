package domain

import (
	"context"

	"gorm.io/datatypes"
)

type SetlistEntry struct {
	ID          int64          `json:"id"`
	SongID      int64          `json:"song_id"`
	SetlistID   int64          `json:"setlist_id"`
	Transpose   int16          `json:"transpose"`
	Notes       string         `json:"notes"`
	Arrangement datatypes.JSON `json:"arrangement"`
}

type SetlistEntryService interface {
	AuthMultiStorer[SetlistEntry]
	Fetcher[SetlistEntry]
	FetchBySetlist(ctx context.Context, setlist *Setlist) (*[]SetlistEntry, error)
	AuthMultiUpdater[SetlistEntry]
	RemoveBatch(ctx context.Context, setlist *Setlist, ids []int64, principal *User) error
	RemoveBySetlist(ctx context.Context, setlist *Setlist, principal *User) error
}

type SetlistEntryRepository interface {
	Getter[SetlistEntry]
	GetBySetlist(ctx context.Context, setlist *Setlist) (*[]SetlistEntry, error)
	Creator[SetlistEntry]
	Updater[SetlistEntry]
	Deleter[SetlistEntry]
}
