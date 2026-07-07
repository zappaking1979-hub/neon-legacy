package application

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/neonlegacy/server/internal/domain/crime"
	"github.com/neonlegacy/server/internal/domain/player"
)

type crimePlayer struct {
	nerve     int
	strength  int
	level     int
	exp       int
	expMax    int
	cash      int64
}

func (p *crimePlayer) Nerve() int    { return p.nerve }
func (p *crimePlayer) Strength() int { return p.strength }
func (p *crimePlayer) Level() int    { return p.level }

type CrimeService struct {
	crimeRepo  crime.Repository
	playerRepo player.Repository
	rng        *rand.Rand
}

func NewCrimeService(crimeRepo crime.Repository, playerRepo player.Repository) *CrimeService {
	return &CrimeService{
		crimeRepo:  crimeRepo,
		playerRepo: playerRepo,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *CrimeService) ListCrimes(ctx context.Context) ([]*crime.Crime, error) {
	return s.crimeRepo.List(ctx)
}

func (s *CrimeService) DoCrime(ctx context.Context, p *player.Player, crimeID, multiplier int) (*crime.Result, error) {
	if multiplier < 1 || multiplier > 10 {
		multiplier = 1
	}

	c, err := s.crimeRepo.GetByID(ctx, crimeID)
	if err != nil {
		return nil, fmt.Errorf("get crime: %w", err)
	}

	if ok, msg := crime.CanCommit(c, p.Level, p.Strength, p.Defense, p.Speed, p.Nerve); !ok {
		return nil, fmt.Errorf("cannot commit: %s", msg)
	}

	totalNerve := c.NerveCost * multiplier
	if p.Nerve < totalNerve {
		return nil, fmt.Errorf("not enough nerve: need %d, have %d", totalNerve, p.Nerve)
	}

	var lastResult *crime.Result
	totalExp := 0
	totalCash := int64(0)

	for i := 0; i < multiplier; i++ {
		cp := &crimePlayer{
			nerve:    p.Nerve,
			strength: p.Strength,
			level:    p.Level,
		}

		result := crime.Commit(c, cp, s.rng)
		lastResult = result

		if result.Jailed {
			p.Nerve -= c.NerveCost
			p.JailTime = time.Now().Add(5 * time.Minute)
			s.playerRepo.Update(ctx, p)
			return result, nil
		}

		if result.Success {
			totalExp += result.ExpGain
			totalCash += int64(result.CashGain)
		}

		p.Nerve -= c.NerveCost
	}

	p.Cash += totalCash
	p.AddExp(totalExp)
	s.playerRepo.Update(ctx, p)

	lastResult.ExpGain = totalExp
	lastResult.CashGain = int(totalCash)
	lastResult.Message = fmt.Sprintf("Completed %dx %s. Earned $%d and %d EXP.",
		multiplier, c.Name, totalCash, totalExp)

	return lastResult, nil
}

func (s *CrimeService) GetCrime(ctx context.Context, id int) (*crime.Crime, error) {
	return s.crimeRepo.GetByID(ctx, id)
}
