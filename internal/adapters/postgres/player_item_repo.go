package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neonlegacy/server/internal/domain/item"
)

const playerItemColumns = `SELECT pi.id, pi.player_id, pi.quantity,
	i.id, i.name, i.description, i.type, i.buy_price, i.sell_price
	FROM player_items pi JOIN items i ON i.id = pi.item_id`

type PlayerItemRepo struct {
	pool *pgxpool.Pool
}

func NewPlayerItemRepo(pool *pgxpool.Pool) *PlayerItemRepo {
	return &PlayerItemRepo{pool: pool}
}

func scanPlayerItem(row interface{ Scan(dest ...any) error }) (*item.PlayerItem, error) {
	var pi item.PlayerItem
	err := row.Scan(
		&pi.ID, &pi.PlayerID, &pi.Quantity,
		&pi.Item.ID, &pi.Item.Name, &pi.Item.Description, &pi.Item.Type, &pi.Item.BuyPrice, &pi.Item.SellPrice,
	)
	if err != nil {
		return nil, err
	}
	return &pi, nil
}

func (r *PlayerItemRepo) ListByPlayer(ctx context.Context, playerID uuid.UUID) ([]item.PlayerItem, error) {
	rows, err := r.pool.Query(ctx, playerItemColumns+" WHERE pi.player_id = $1 ORDER BY i.id", playerID)
	if err != nil {
		return nil, fmt.Errorf("list player items: %w", err)
	}
	defer rows.Close()

	var results []item.PlayerItem
	for rows.Next() {
		pi, err := scanPlayerItem(rows)
		if err != nil {
			return nil, fmt.Errorf("scan player item: %w", err)
		}
		results = append(results, *pi)
	}
	return results, rows.Err()
}

func (r *PlayerItemRepo) GetByPlayerAndItem(ctx context.Context, playerID uuid.UUID, itemID int) (*item.PlayerItem, error) {
	row := r.pool.QueryRow(ctx, playerItemColumns+" WHERE pi.player_id = $1 AND pi.item_id = $2", playerID, itemID)
	pi, err := scanPlayerItem(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get player item: %w", err)
	}
	return pi, nil
}

func (r *PlayerItemRepo) Add(ctx context.Context, playerID uuid.UUID, itemID int, quantity int) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO player_items (player_id, item_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (player_id, item_id) DO UPDATE SET quantity = player_items.quantity + $3`,
		playerID, itemID, quantity)
	if err != nil {
		return fmt.Errorf("add player item: %w", err)
	}
	return nil
}

func (r *PlayerItemRepo) Remove(ctx context.Context, playerID uuid.UUID, itemID int, quantity int) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE player_items SET quantity = quantity - $3
		WHERE player_id = $1 AND item_id = $2 AND quantity >= $3`,
		playerID, itemID, quantity)
	if err != nil {
		return fmt.Errorf("remove player item: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("not enough items")
	}

	r.pool.Exec(ctx, `DELETE FROM player_items WHERE player_id = $1 AND item_id = $2 AND quantity <= 0`, playerID, itemID)
	return nil
}

func (r *PlayerItemRepo) SetQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	_, err := r.pool.Exec(ctx, `UPDATE player_items SET quantity = $2 WHERE id = $1`, id, quantity)
	if err != nil {
		return fmt.Errorf("set quantity: %w", err)
	}
	return nil
}
