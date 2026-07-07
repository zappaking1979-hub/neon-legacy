package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neonlegacy/server/internal/domain/item"
)

const itemColumns = "SELECT id, name, description, type, buy_price, sell_price FROM items"

type ItemRepo struct {
	pool *pgxpool.Pool
}

func NewItemRepo(pool *pgxpool.Pool) *ItemRepo {
	return &ItemRepo{pool: pool}
}

func (r *ItemRepo) List(ctx context.Context) ([]item.Item, error) {
	rows, err := r.pool.Query(ctx, itemColumns+" ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	defer rows.Close()

	var results []item.Item
	for rows.Next() {
		var it item.Item
		if err := rows.Scan(&it.ID, &it.Name, &it.Description, &it.Type, &it.BuyPrice, &it.SellPrice); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		results = append(results, it)
	}
	return results, rows.Err()
}

func (r *ItemRepo) GetByID(ctx context.Context, id int) (*item.Item, error) {
	row := r.pool.QueryRow(ctx, itemColumns+" WHERE id = $1", id)
	var it item.Item
	if err := row.Scan(&it.ID, &it.Name, &it.Description, &it.Type, &it.BuyPrice, &it.SellPrice); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get item: %w", err)
	}
	return &it, nil
}
