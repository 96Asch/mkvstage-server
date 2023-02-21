package domain

type UserRole struct {
	ID  int64 `json:"id"`
	UID int64 `json:"uid"`
	RID int64 `json:"rid"`
}

type UserRoleService interface {
	AuthMultiUpdater[UserRole]
}

type UserRoleRepository interface {
	Creator[UserRole]
	Getter[UserRole]
	Updater[UserRole]
	Deleter[UserRole]
}
