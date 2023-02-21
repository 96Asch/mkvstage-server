package domain

import "context"

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" gorm:"type:varchar(255);unique"`
	Description string `json:"description"`
}

type RoleService interface {
	AuthSingleStorer[Role]
	Fetcher[Role]
	AuthSingleUpdater[Role]
	AuthSingleRemover[Role]
}

type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	Getter[Role]
	Update(ctx context.Context, role *Role) error
	Deleter(ctx context.Context, rid int64) error
}
