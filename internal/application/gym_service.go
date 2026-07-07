package application

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/neonlegacy/server/internal/domain/gym"
	"github.com/neonlegacy/server/internal/domain/player"
)

type GymService struct {
	gymRepo    gym.Repository
	playerRepo player.Repository
	rng        *rand.Rand
}

func NewGymService(gr gym.Repository, pr player.Repository) *GymService {
	return &GymService{
		gymRepo:    gr,
		playerRepo: pr,
		rng:        rand.New(rand.NewSource(rand.Int63())),
	}
}

func (s *GymService) ListExercises(ctx context.Context) ([]gym.Exercise, error) {
	return s.gymRepo.List(ctx)
}

func (s *GymService) Train(ctx context.Context, p *player.Player, exerciseID int) (*gym.Result, error) {
	ex, err := s.gymRepo.GetByID(ctx, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("get exercise: %w", err)
	}
	if ex == nil {
		return nil, fmt.Errorf("exercise not found")
	}

	ok, reason := ex.CanTrain(p.Level, p.Energy)
	if !ok {
		return &gym.Result{
			Success:  false,
			Message:  reason,
			StatGain: 0,
		}, nil
	}

	p.SpendEnergy(ex.EnergyCost)

	gain := ex.GainMin + s.rng.Intn(ex.GainMax-ex.GainMin+1)

	switch ex.Stat {
	case gym.StatStrength:
		p.TrainStrength(gain)
	case gym.StatDefense:
		p.TrainDefense(gain)
	case gym.StatSpeed:
		p.TrainSpeed(gain)
	case gym.StatAgility:
		p.TrainAgility(gain)
	}

	expGain := gain * 5
	p.AddExp(expGain)

	if err := s.playerRepo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update player: %w", err)
	}

	return &gym.Result{
		Success:  true,
		StatGain: gain,
		StatName: string(ex.Stat),
		Message:  ex.Description,
		ExpGain:  expGain,
	}, nil
}

func (s *GymService) GetExercise(ctx context.Context, id int) (*gym.Exercise, error) {
	return s.gymRepo.GetByID(ctx, id)
}
