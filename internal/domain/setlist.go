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
	Global    bool           `json:"is_global"`
	Deadline  time.Time      `json:"deadline"`
	UpdatedAt time.Time      `json:"updated_at"`
	Order     datatypes.JSON `json:"-"`
}

type SetlistService interface {
	AuthSingleStorer[Setlist]
	Fetcher[Setlist]
	FetchAllGlobal(ctx context.Context, principal *User) (*[]Setlist, error)
	Update(ctx context.Context, setlist *Setlist, principal *User) (*Setlist, error)
	AuthSingleRemover[Setlist]
}

type SetlistRepository interface {
	Getter[Setlist]
	GetAllGlobal(ctx context.Context, uid int64) (*[]Setlist, error)
	Delete(ctx context.Context, slid int64) error
	Update(ctx context.Context, setlist *Setlist) (*Setlist, error)
	Create(ctx context.Context, setlist *Setlist) error
}
