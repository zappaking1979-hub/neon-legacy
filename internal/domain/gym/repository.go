package gym

import "context"

type Repository interface {
	List(ctx context.Context) ([]Exercise, error)
	GetByID(ctx context.Context, id int) (*Exercise, error)
}
