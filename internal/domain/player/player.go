package player

import (
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type Player struct {
	ID           uuid.UUID
	Email        string
	Username     string
	PasswordHash string
	Gender       Gender
	CreatedAt    time.Time
	LastActiveAt time.Time

	Level    int
	Prestige int
	Exp      int
	ExpMax   int

	HP      int
	HPMax   int
	Energy  int
	EnergyMax int
	Nerve   int
	NerveMax int
	Awake   int
	AwakeMax int

	Strength int
	Defense  int
	Speed    int
	Agility  int

	Cash    int64
	Bank    int64
	Points  int64
	Credits int64

	GangID *uuid.UUID

	HospitalTime time.Time
	JailTime     time.Time
}

func NewPlayer(email, username, passwordHash string, gender Gender) *Player {
	now := time.Now()
	return &Player{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		Gender:       gender,
		CreatedAt:    now,
		LastActiveAt: now,

		Level:    1,
		Prestige: 0,
		Exp:      0,
		ExpMax:   100,

		HP:      100,
		HPMax:   100,
		Energy:  100,
		EnergyMax: 100,
		Nerve:   50,
		NerveMax: 50,
		Awake:   100,
		AwakeMax: 100,

		Strength: 10,
		Defense:  10,
		Speed:    10,
		Agility:  10,

		Cash:    0,
		Bank:    0,
		Points:  0,
		Credits: 0,
	}
}

func (p *Player) IsInHospital() bool {
	return time.Now().Before(p.HospitalTime)
}

func (p *Player) IsInJail() bool {
	return time.Now().Before(p.JailTime)
}

func (p *Player) CanAct() bool {
	return !p.IsInHospital() && !p.IsInJail()
}

func (p *Player) AddExp(amount int) {
	p.Exp += amount
	for p.Exp >= p.ExpMax && p.ExpMax > 0 {
		p.Exp -= p.ExpMax
		p.LevelUp()
	}
}

func (p *Player) LevelUp() {
	p.Level++
	p.ExpMax = expForLevel(p.Level)
	p.HPMax = 100 + (p.Level-1)*10
	p.EnergyMax = 100 + (p.Level-1)*2
	p.NerveMax = 50 + (p.Level-1)
	p.HP = p.HPMax
	p.Energy = p.EnergyMax
	p.Nerve = p.NerveMax
}

func expForLevel(level int) int {
	return 100 * level * level
}

func (p *Player) RegenHP(amount int) {
	p.HP += amount
	if p.HP > p.HPMax {
		p.HP = p.HPMax
	}
}

func (p *Player) RegenEnergy(amount int) {
	p.Energy += amount
	if p.Energy > p.EnergyMax {
		p.Energy = p.EnergyMax
	}
}

func (p *Player) RegenNerve(amount int) {
	p.Nerve += amount
	if p.Nerve > p.NerveMax {
		p.Nerve = p.NerveMax
	}
}

func (p *Player) RegenAwake(amount int) {
	p.Awake += amount
	if p.Awake > p.AwakeMax {
		p.Awake = p.AwakeMax
	}
}
