package crime

import (
	"math/rand"
)

type Crime struct {
	ID           int
	Name         string
	Description  string
	NerveCost    int
	MinLevel     int
	MinStrength  int
	MinDefense   int
	MinSpeed     int
	SuccessRate  float64
	JailChance   float64
	ExpReward    int
	CashRewardMin int
	CashRewardMax int
}

type Result struct {
	Success  bool
	Jailed   bool
	ExpGain  int
	CashGain int
	Message  string
}

func CanCommit(crime *Crime, level, strength, defense, speed, nerve int) (bool, string) {
	if nerve < crime.NerveCost {
		return false, "not enough nerve"
	}
	if level < crime.MinLevel {
		return false, "level too low"
	}
	if strength < crime.MinStrength {
		return false, "strength too low"
	}
	if defense < crime.MinDefense {
		return false, "defense too low"
	}
	if speed < crime.MinSpeed {
		return false, "speed too low"
	}
	return true, ""
}

func Commit(crime *Crime, playerStats interface {
	Nerve() int
	Strength() int
	Level() int
}, rng *rand.Rand) *Result {
	roll := rng.Float64() * 100

	if roll < crime.JailChance {
		return &Result{
			Jailed:  true,
			Message: "You got caught! The cops threw you in jail.",
		}
	}

	if roll < crime.JailChance+crime.SuccessRate {
		cash := crime.CashRewardMin + rng.Intn(crime.CashRewardMax-crime.CashRewardMin+1)
		exp := crime.ExpReward + rng.Intn(crime.ExpReward)
		return &Result{
			Success:  true,
			ExpGain:  exp,
			CashGain: cash,
			Message:  "Success! You pulled it off clean.",
		}
	}

	return &Result{
		Success: false,
		Message: "You failed. No reward, but at least you are not in jail.",
	}
}

func ExperienceForLevel(level int) int {
	return 100 * level * level
}

func NerveRegenPerMinute() int {
	return 1
}

func EnergyRegenPerMinute() int {
	return 2
}

func HPRegenPerMinute() int {
	return 5
}
