package player

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	p := NewPlayer("test@test.com", "testuser", "hash", GenderMale)

	if p.Email != "test@test.com" {
		t.Errorf("expected test@test.com, got %s", p.Email)
	}
	if p.Level != 1 {
		t.Errorf("expected level 1, got %d", p.Level)
	}
	if p.HPMax != 100 {
		t.Errorf("expected HP max 100, got %d", p.HPMax)
	}
	if p.Strength != 10 {
		t.Errorf("expected strength 10, got %d", p.Strength)
	}
	if p.Cash != 0 {
		t.Errorf("expected cash 0, got %d", p.Cash)
	}
}

func TestAddExp(t *testing.T) {
	p := NewPlayer("test@test.com", "testuser", "hash", GenderMale)

	if p.Level != 1 {
		t.Fatalf("expected level 1, got %d", p.Level)
	}

	p.AddExp(50)
	if p.Level != 1 {
		t.Errorf("should not level up with 50 exp, got level %d", p.Level)
	}
	if p.Exp != 50 {
		t.Errorf("expected exp 50, got %d", p.Exp)
	}

	p.AddExp(100)
	if p.Level != 2 {
		t.Errorf("expected level 2 after 150 total exp, got %d", p.Level)
	}
	if p.Exp != 50 {
		t.Errorf("expected 50 overflow exp, got %d", p.Exp)
	}
}

func TestLevelUpStats(t *testing.T) {
	p := NewPlayer("test@test.com", "testuser", "hash", GenderMale)
	p.LevelUp()

	if p.Level != 2 {
		t.Errorf("expected level 2, got %d", p.Level)
	}
	if p.HPMax != 110 {
		t.Errorf("expected HP max 110, got %d", p.HPMax)
	}
	if p.EnergyMax != 102 {
		t.Errorf("expected energy max 102, got %d", p.EnergyMax)
	}
	if p.NerveMax != 51 {
		t.Errorf("expected nerve max 51, got %d", p.NerveMax)
	}
	if p.HP != p.HPMax {
		t.Errorf("expected HP to be refilled to max on level up")
	}
}

func TestRegen(t *testing.T) {
	p := NewPlayer("test@test.com", "testuser", "hash", GenderMale)
	p.HP = 50
	p.Energy = 50
	p.Nerve = 25
	p.Awake = 50

	p.RegenHP(10)
	if p.HP != 60 {
		t.Errorf("expected HP 60, got %d", p.HP)
	}

	p.RegenHP(100)
	if p.HP > p.HPMax {
		t.Errorf("HP should not exceed max, got %d", p.HP)
	}

	p.RegenEnergy(10)
	if p.Energy != 60 {
		t.Errorf("expected energy 60, got %d", p.Energy)
	}

	p.RegenNerve(10)
	if p.Nerve != 35 {
		t.Errorf("expected nerve 35, got %d", p.Nerve)
	}

	p.RegenAwake(10)
	if p.Awake != 60 {
		t.Errorf("expected awake 60, got %d", p.Awake)
	}
}

func TestJailAndHospital(t *testing.T) {
	p := NewPlayer("test@test.com", "testuser", "hash", GenderMale)

	if !p.CanAct() {
		t.Error("new player should be able to act")
	}

	// Jail and hospital times are in the past by default, so CanAct should be true
	if p.IsInJail() {
		t.Error("new player should not be in jail")
	}
	if p.IsInHospital() {
		t.Error("new player should not be in hospital")
	}
}

func TestExpForLevel(t *testing.T) {
	tests := []struct {
		level int
		exp   int
	}{
		{1, 100},
		{2, 400},
		{5, 2500},
		{10, 10000},
	}

	for _, tt := range tests {
		got := expForLevel(tt.level)
		if got != tt.exp {
			t.Errorf("expForLevel(%d) = %d, want %d", tt.level, got, tt.exp)
		}
	}
}

func TestMultipleLevelUps(t *testing.T) {
	p := NewPlayer("test@test.com", "testuser", "hash", GenderMale)

	p.AddExp(10000)

	if p.Level < 5 {
		t.Errorf("expected at least level 5 after 10000 exp, got %d", p.Level)
	}
	if p.HP != p.HPMax {
		t.Errorf("HP should be refilled after multi-level-up")
	}
}

func TestGenderValidation(t *testing.T) {
	tests := map[Gender]bool{
		GenderMale:   true,
		GenderFemale: true,
		GenderOther:  true,
		Gender(""):   false,
		Gender("x"):  false,
	}

	for g, valid := range tests {
		switch g {
		case GenderMale, GenderFemale, GenderOther:
			if !valid {
				t.Errorf("expected %q to be valid", g)
			}
		default:
			if valid {
				t.Errorf("expected %q to be invalid", g)
			}
		}
	}
}
