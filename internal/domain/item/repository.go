package item

import (
	"context"

	"github.com/google/uuid"
)

type ItemRepository interface {
	List(ctx context.Context) ([]Item, error)
	GetByID(ctx context.Context, id int) (*Item, error)
}

type PlayerItemRepository interface {
	ListByPlayer(ctx context.Context, playerID uuid.UUID) ([]PlayerItem, error)
	GetByPlayerAndItem(ctx context.Context, playerID uuid.UUID, itemID int) (*PlayerItem, error)
	Add(ctx context.Context, playerID uuid.UUID, itemID int, quantity int) error
	Remove(ctx context.Context, playerID uuid.UUID, itemID int, quantity int) error
	SetQuantity(ctx context.Context, id uuid.UUID, quantity int) error
}
