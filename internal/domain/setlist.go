package domain

import (
	"context"
	"time"

	"gorm.io/datatypes"
)

type Setlist struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	CreatorID int64          `json:"creator_id"`
	Deadline  time.Time      `json:"deadline"`
	UpdatedAt time.Time      `json:"updated_at"`
	Order     datatypes.JSON `json:"-"`
}

type SetlistService interface {
	AuthSingleStorer[Setlist]
	Fetcher[Setlist]
	FetchByTimeframe(ctx context.Context, from time.Time, to time.Time) (*[]Setlist, error)
	Update(ctx context.Context, setlist *Setlist, principal *User) (*Setlist, error)
	AuthSingleRemover[Setlist]
}

type SetlistRepository interface {
	Getter[Setlist]
	GetByTimeframe(ctx context.Context, from time.Time, to time.Time) (*[]Setlist, error)
	Delete(ctx context.Context, slid int64) error
	Update(ctx context.Context, setlist *Setlist) (*Setlist, error)
	Create(ctx context.Context, setlist *Setlist) error
}
