package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neonlegacy/server/internal/domain/player"
)

type PlayerRepo struct {
	pool *pgxpool.Pool
}

func NewPlayerRepo(pool *pgxpool.Pool) *PlayerRepo {
	return &PlayerRepo{pool: pool}
}

func (r *PlayerRepo) Create(ctx context.Context, p *player.Player) error {
	query := `
		INSERT INTO players (
			id, email, username, password_hash, gender,
			created_at, last_active_at,
			level, prestige, exp, exp_max,
			hp, hp_max, energy, energy_max, nerve, nerve_max, awake, awake_max,
			strength, defense, speed, agility,
			cash, bank, points, credits,
			hospital_time, jail_time
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7,
			$8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18, $19,
			$20, $21, $22, $23,
			$24, $25, $26, $27,
			$28, $29
		)`

	_, err := r.pool.Exec(ctx, query,
		p.ID, p.Email, p.Username, p.PasswordHash, p.Gender,
		p.CreatedAt, p.LastActiveAt,
		p.Level, p.Prestige, p.Exp, p.ExpMax,
		p.HP, p.HPMax, p.Energy, p.EnergyMax, p.Nerve, p.NerveMax, p.Awake, p.AwakeMax,
		p.Strength, p.Defense, p.Speed, p.Agility,
		p.Cash, p.Bank, p.Points, p.Credits,
		p.HospitalTime, p.JailTime,
	)
	if err != nil {
		return fmt.Errorf("insert player: %w", err)
	}
	return nil
}

func (r *PlayerRepo) GetByID(ctx context.Context, id uuid.UUID) (*player.Player, error) {
	return r.scanOne(ctx, r.pool.QueryRow(ctx, "SELECT * FROM players WHERE id = $1", id))
}

func (r *PlayerRepo) GetByEmail(ctx context.Context, email string) (*player.Player, error) {
	return r.scanOne(ctx, r.pool.QueryRow(ctx, "SELECT * FROM players WHERE email = $1", email))
}

func (r *PlayerRepo) GetByUsername(ctx context.Context, username string) (*player.Player, error) {
	return r.scanOne(ctx, r.pool.QueryRow(ctx, "SELECT * FROM players WHERE username = $1", username))
}

func (r *PlayerRepo) Update(ctx context.Context, p *player.Player) error {
	query := `
		UPDATE players SET
			email = $2, password_hash = $3, last_active_at = $4,
			level = $5, prestige = $6, exp = $7, exp_max = $8,
			hp = $9, hp_max = $10, energy = $11, energy_max = $12,
			nerve = $13, nerve_max = $14, awake = $15, awake_max = $16,
			strength = $17, defense = $18, speed = $19, agility = $20,
			cash = $21, bank = $22, points = $23, credits = $24,
			hospital_time = $25, jail_time = $26,
			gang_id = $27
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query,
		p.ID, p.Email, p.PasswordHash, p.LastActiveAt,
		p.Level, p.Prestige, p.Exp, p.ExpMax,
		p.HP, p.HPMax, p.Energy, p.EnergyMax,
		p.Nerve, p.NerveMax, p.Awake, p.AwakeMax,
		p.Strength, p.Defense, p.Speed, p.Agility,
		p.Cash, p.Bank, p.Points, p.Credits,
		p.HospitalTime, p.JailTime,
		p.GangID,
	)
	if err != nil {
		return fmt.Errorf("update player: %w", err)
	}
	return nil
}

func (r *PlayerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM players WHERE id = $1", id)
	return err
}

func (r *PlayerRepo) Search(ctx context.Context, query string, limit, offset int) ([]*player.Player, error) {
	sql := "SELECT * FROM players WHERE username ILIKE $1 ORDER BY level DESC LIMIT $2 OFFSET $3"
	rows, err := r.pool.Query(ctx, sql, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*player.Player
	for rows.Next() {
		p, err := scanPlayer(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, rows.Err()
}

func (r *PlayerRepo) scanOne(ctx context.Context, row pgx.Row) (*player.Player, error) {
	p, err := scanPlayer(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan player: %w", err)
	}
	return p, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanPlayer(row scanner) (*player.Player, error) {
	p := &player.Player{}
	err := row.Scan(
		&p.ID, &p.Email, &p.Username, &p.PasswordHash, &p.Gender,
		&p.CreatedAt, &p.LastActiveAt,
		&p.Level, &p.Prestige, &p.Exp, &p.ExpMax,
		&p.HP, &p.HPMax, &p.Energy, &p.EnergyMax,
		&p.Nerve, &p.NerveMax, &p.Awake, &p.AwakeMax,
		&p.Strength, &p.Defense, &p.Speed, &p.Agility,
		&p.Cash, &p.Bank, &p.Points, &p.Credits,
		&p.HospitalTime, &p.JailTime,
		&p.GangID,
	)
	return p, err
}
