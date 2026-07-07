package crime

import "context"

type Repository interface {
	List(ctx context.Context) ([]*Crime, error)
	GetByID(ctx context.Context, id int) (*Crime, error)
}
