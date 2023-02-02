package domain

import "context"

type Fetcher[T any] interface {
	FetchByID(ctx context.Context, id int64) (*T, error)
	FetchAll(ctx context.Context) ([]*T, error)
}
