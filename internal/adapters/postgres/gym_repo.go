package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neonlegacy/server/internal/domain/gym"
)

const gymColumns = "SELECT id, name, description, stat, energy_cost, min_level, gain_min, gain_max FROM gym_exercises"

type GymRepo struct {
	pool *pgxpool.Pool
}

func NewGymRepo(pool *pgxpool.Pool) *GymRepo {
	return &GymRepo{pool: pool}
}

func (r *GymRepo) List(ctx context.Context) ([]gym.Exercise, error) {
	rows, err := r.pool.Query(ctx, gymColumns+" ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("list gym exercises: %w", err)
	}
	defer rows.Close()

	var results []gym.Exercise
	for rows.Next() {
		var e gym.Exercise
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Stat, &e.EnergyCost, &e.MinLevel, &e.GainMin, &e.GainMax); err != nil {
			return nil, fmt.Errorf("scan gym exercise: %w", err)
		}
		results = append(results, e)
	}
	return results, rows.Err()
}

func (r *GymRepo) GetByID(ctx context.Context, id int) (*gym.Exercise, error) {
	row := r.pool.QueryRow(ctx, gymColumns+" WHERE id = $1", id)
	var e gym.Exercise
	if err := row.Scan(&e.ID, &e.Name, &e.Description, &e.Stat, &e.EnergyCost, &e.MinLevel, &e.GainMin, &e.GainMax); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get gym exercise: %w", err)
	}
	return &e, nil
}
