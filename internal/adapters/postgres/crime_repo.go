package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neonlegacy/server/internal/domain/crime"
)

type CrimeRepo struct {
	pool *pgxpool.Pool
}

func NewCrimeRepo(pool *pgxpool.Pool) *CrimeRepo {
	return &CrimeRepo{pool: pool}
}

func (r *CrimeRepo) List(ctx context.Context) ([]*crime.Crime, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, description, nerve_cost, min_level,
		       min_strength, min_defense, min_speed,
		       success_rate, jail_chance, exp_reward,
		       cash_reward_min, cash_reward_max
		FROM crimes ORDER BY nerve_cost ASC`)
	if err != nil {
		return nil, fmt.Errorf("query crimes: %w", err)
	}
	defer rows.Close()

	var crimes []*crime.Crime
	for rows.Next() {
		c := &crime.Crime{}
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.NerveCost, &c.MinLevel,
			&c.MinStrength, &c.MinDefense, &c.MinSpeed,
			&c.SuccessRate, &c.JailChance, &c.ExpReward,
			&c.CashRewardMin, &c.CashRewardMax,
		); err != nil {
			return nil, fmt.Errorf("scan crime: %w", err)
		}
		crimes = append(crimes, c)
	}
	return crimes, rows.Err()
}

func (r *CrimeRepo) GetByID(ctx context.Context, id int) (*crime.Crime, error) {
	c := &crime.Crime{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, description, nerve_cost, min_level,
		       min_strength, min_defense, min_speed,
		       success_rate, jail_chance, exp_reward,
		       cash_reward_min, cash_reward_max
		FROM crimes WHERE id = $1`, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.NerveCost, &c.MinLevel,
		&c.MinStrength, &c.MinDefense, &c.MinSpeed,
		&c.SuccessRate, &c.JailChance, &c.ExpReward,
		&c.CashRewardMin, &c.CashRewardMax,
	)
	if err != nil {
		return nil, fmt.Errorf("get crime %d: %w", id, err)
	}
	return c, nil
}
