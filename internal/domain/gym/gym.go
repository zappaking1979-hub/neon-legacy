package gym

import "strconv"

type Stat string

const (
	StatStrength Stat = "strength"
	StatDefense  Stat = "defense"
	StatSpeed    Stat = "speed"
	StatAgility  Stat = "agility"
)

func ValidStat(s Stat) bool {
	switch s {
	case StatStrength, StatDefense, StatSpeed, StatAgility:
		return true
	}
	return false
}

type Exercise struct {
	ID         int
	Name       string
	Description string
	Stat       Stat
	EnergyCost int
	MinLevel   int
	GainMin    int
	GainMax    int
}

func (e *Exercise) CanTrain(playerLevel int, playerEnergy int) (bool, string) {
	if playerLevel < e.MinLevel {
		return false, "Level " + strconv.Itoa(e.MinLevel) + " required"
	}
	if playerEnergy < e.EnergyCost {
		return false, "Not enough energy"
	}
	return true, ""
}

type Result struct {
	Success   bool
	StatGain  int
	StatName  string
	Message   string
	ExpGain   int
}
