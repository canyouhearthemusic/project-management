package project

import "context"

type Repository interface {
	Create(context.Context, Entity) (string, Entity, error)
	Search(ctx context.Context, filter, value string) ([]Entity, error)
	List(ctx context.Context) ([]Entity, error)
	Get(ctx context.Context, id string) (Entity, error)
	Update(ctx context.Context, id string, p Entity) error
	Delete(ctx context.Context, id string) error
}
