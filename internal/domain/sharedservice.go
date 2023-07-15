package domain

import "context"

type Fetcher[T any] interface {
	FetchByID(ctx context.Context, id int64) (*T, error)
	FetchAll(ctx context.Context) (*[]T, error)
}

type AuthSingleUpdater[T any] interface {
	Update(ctx context.Context, domain *T, principal *User) error
}

type AuthMultiUpdater[T any] interface {
	UpdateBatch(ctx context.Context, domain *[]T, principal *User) error
}

type AuthUpdater[T any] interface {
	AuthSingleUpdater[T]
	AuthMultiUpdater[T]
}

type AuthSingleStorer[T any] interface {
	Store(ctx context.Context, domain *T, principal *User) error
}

type AuthMultiStorer[T any] interface {
	StoreBatch(ctx context.Context, domain *[]T, principal *User) error
}

type AuthSingleRemover[T any] interface {
	Remove(ctx context.Context, id int64, principal *User) error
}

type AuthMultiRemover[T any] interface {
	RemoveBatch(ctx context.Context, ids []int64, principal *User) error
}
