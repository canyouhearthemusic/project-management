package task

import "context"

type Repository interface {
	List(ctx context.Context) ([]Entity, error)
	Search(ctx context.Context, filter, value string) ([]Entity, error)
	Get(ctx context.Context, id string) (Entity, error)
	Create(ctx context.Context, Entity Entity) (string, Entity, error)
	Update(ctx context.Context, id string, Entity Entity) error
	Delete(ctx context.Context, id string) error
}
