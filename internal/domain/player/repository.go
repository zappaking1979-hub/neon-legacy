package player

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, p *Player) error
	GetByID(ctx context.Context, id uuid.UUID) (*Player, error)
	GetByEmail(ctx context.Context, email string) (*Player, error)
	GetByUsername(ctx context.Context, username string) (*Player, error)
	Update(ctx context.Context, p *Player) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string, limit, offset int) ([]*Player, error)
}
