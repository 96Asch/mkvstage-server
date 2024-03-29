package domain

import (
	"context"
	"time"
)

type Setlist struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatorID int64     `json:"creator_id"`
	Deadline  time.Time `json:"deadline"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SetlistService interface {
	AuthSingleStorer[Setlist]
	Fetcher[Setlist]
	Fetch(ctx context.Context, from time.Time, to time.Time) (*[]Setlist, error)
	FetchByTimeframe(ctx context.Context, from time.Time, to time.Time) (*[]Setlist, error)
	Update(ctx context.Context, setlist *Setlist, principal *User) (*Setlist, error)
	AuthSingleRemover[Setlist]
}

type SetlistRepository interface {
	GetByID(ctx context.Context, id int64) (*Setlist, error)
	GetByIDs(ctx context.Context, id []int64) (*[]Setlist, error)
	GetAll(ctx context.Context) (*[]Setlist, error)
	Get(ctx context.Context, from time.Time, to time.Time) (*[]Setlist, error)
	GetByTimeframe(ctx context.Context, from time.Time, to time.Time) (*[]Setlist, error)
	Delete(ctx context.Context, slid int64) error
	Update(ctx context.Context, setlist *Setlist) (*Setlist, error)
	Create(ctx context.Context, setlist *Setlist) error
}
