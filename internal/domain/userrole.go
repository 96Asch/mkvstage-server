package domain

import "context"

type UserRole struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"uid" gorm:"uniqueIndex:user_role"`
	RoleID int64 `json:"rid" gorm:"uniqueIndex:user_role"`
	Active bool  `json:"active"`
}

type UserRoleService interface {
	AuthMultiUpdater[UserRole]
}

type UserRoleRepository interface {
	Creator[UserRole]
	Getter[UserRole]
	Updater[UserRole]
	Deleter[UserRole]
	DeleteByRID(ctx context.Context, rid int64) error
}
