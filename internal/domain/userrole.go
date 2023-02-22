package domain

import "context"

type UserRole struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"uid" gorm:"uniqueIndex:user_role"`
	RoleID int64 `json:"rid" gorm:"uniqueIndex:user_role"`
	Active bool  `json:"active"`
}

type UserRoleService interface {
	GetAll(ctx context.Context) (*[]UserRole, error)
	GetByUser(ctx context.Context, user *User) (*[]UserRole, error)
	SetActiveBatch(ctx context.Context, urids []int64, principal *User) (*[]UserRole, error)
}

type UserRoleRepository interface {
	Creator[UserRole]
	Getter[UserRole]
	GetByUID(ctx context.Context, uid int64) (*[]UserRole, error)
	Updater[UserRole]
	Deleter[UserRole]
	DeleteByRID(ctx context.Context, rid int64) error
	DeleteByUID(ctx context.Context, uid int64) error
}
