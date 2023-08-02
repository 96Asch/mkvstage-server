package domain

import (
	"context"

	"gorm.io/gorm"
)

type UserRole struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"-" gorm:"uniqueIndex:user_role"`
	User      *User          `json:"user"`
	RoleID    int64          `json:"-" gorm:"uniqueIndex:user_role"`
	Role      *Role          `json:"role"`
	Active    bool           `json:"active"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type UserRoleService interface {
	FetchAll(ctx context.Context) (*[]UserRole, error)
	FetchByUser(ctx context.Context, user *User) (*[]UserRole, error)
	SetActiveBatch(ctx context.Context, urids []int64, principal *User) (*[]UserRole, error)
}

type UserRoleRepository interface {
	Creator[UserRole]
	Getter[UserRole]
	Get(ctx context.Context, ids []int64) (*[]UserRole, error)
	GetByUID(ctx context.Context, uid int64) (*[]UserRole, error)
	Updater[UserRole]
	Deleter[UserRole]
	DeleteByRID(ctx context.Context, rid int64) error
	DeleteByUID(ctx context.Context, uid int64) error
}
