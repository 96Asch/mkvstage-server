package domain

import "context"

type SetlistRole struct {
	ID        int64 `json:"id"`
	SetlistID int64 `json:"setlist_id"`
	UserRole  int64 `json:"userrole_id"`
}

type SetlistRoleService interface {
	Fetch(ctx context.Context, setlists *[]Setlist) (*[]SetlistRole, error)
	Store(ctx context.Context, setlistRoles *[]SetlistRole) error
	Update(ctx context.Context, setlistRoles *[]SetlistRole) error
	Remove(ctx context.Context, setlistRoleIDs []int64) error
}

type SetlistRoleRepository interface {
	Create(ctx context.Context, setlistRoles *[]SetlistRole) error
	Get(ctx context.Context, setlistIDs []int64) (*[]SetlistRole, error)
	Update(ctx context.Context, setlistRoles *[]SetlistRole) error
	Delete(ctx context.Context, setlistRoleIDs []int64) error
}
