package domain

import "context"

type Fetcher[T any] interface {
	FetchByID(ctx context.Context, id int64) (*T, error)
	FetchAll(ctx context.Context) (*[]T, error)
}

type AuthSingleUpdater[T any] interface {
	Update(ctx context.Context, domain *T, principal *User) error
}

type AuthSingleStorer[T any] interface {
	Store(ctx context.Context, domain *T, principal *User) error
}

type AuthSingleRemover[T any] interface {
	Remove(ctx context.Context, id int64, principal *User) error
}
