package user

import "context"

type Repository interface {
	List(ctx context.Context) ([]Entity, error)
	Search(ctx context.Context, filter, value string) ([]Entity, error)
	Create(context.Context, Entity) (string, Entity, error)
	Get(ctx context.Context, id string) (Entity, error)
	Update(ctx context.Context, id string, u Entity) error
	Delete(ctx context.Context, id string) error
}
