package domain

import (
	"context"

	"gorm.io/gorm"
)

type Bundle struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name" gorm:"type:varchar(255);uniqueIndex:name_id" `
	ParentID  int64          `json:"parent_id" gorm:"uniqueIndex:name_id"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type BundleService interface {
	Fetcher[Bundle]
	Store(ctx context.Context, bundle *Bundle, principal *User) error
	Remove(ctx context.Context, bid int64, principal *User) error
	Update(ctx context.Context, bundle *Bundle, principal *User) error
}

type BundleRepository interface {
	Getter[Bundle]
	Create(ctx context.Context, bundle *Bundle) error
	Delete(ctx context.Context, bid int64) error
	Update(ctx context.Context, bundle *Bundle) error
	GetLeaves(ctx context.Context) (*[]Bundle, error)
}
