package domain

import "context"

type Getter[T any] interface {
	GetByID(ctx context.Context, id int64) (*T, error)
	GetAll(ctx context.Context, id int64) ([]*T, error)
}

type Creator[T any] interface {
	Create(ctx context.Context, obj *T) error
	CreateBatch(ctx context.Context, obj []*T) error
}

type Updater[T any] interface {
	Update(ctx context.Context, obj *T) error
	UpdateBatch(ctx context.Context, obj []*T) error
}

type Deleter[T any] interface {
	Delete(ctx context.Context, id int64) error
	DeleteBatch(ctx context.Context, ids []int64) error
}
