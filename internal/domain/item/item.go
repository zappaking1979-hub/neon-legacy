package item

import "github.com/google/uuid"

type Type string

const (
	Consumable Type = "consumable"
	Weapon     Type = "weapon"
	Armor      Type = "armor"
	Special    Type = "special"
)

func ValidType(t Type) bool {
	switch t {
	case Consumable, Weapon, Armor, Special:
		return true
	}
	return false
}

type Item struct {
	ID          int
	Name        string
	Description string
	Type        Type
	BuyPrice    int64
	SellPrice   int64
}

type PlayerItem struct {
	ID       uuid.UUID
	PlayerID uuid.UUID
	Item     Item
	Quantity int
}

func (pi *PlayerItem) Use(p interface {
	SpendEnergy(int) bool
	RegenHP(int)
	RegenEnergy(int)
	RegenNerve(int)
	RegenAwake(int)
}) string {
	switch pi.Item.Type {
	case Consumable:
		switch pi.Item.ID {
		case 1:
			p.RegenHP(25)
			return "You used a First Aid Kit and restored 25 HP."
		case 2:
			p.RegenEnergy(20)
			return "You drank an Energy Drink and restored 20 Energy."
		case 3:
			p.RegenNerve(10)
			return "You took Nerve Pills and restored 10 Nerve."
		case 4:
			p.RegenAwake(15)
			return "You drank Coffee and restored 15 Awake."
		}
	case Special:
		switch pi.Item.ID {
		case 7:
			return "Exp Booster activated. Double EXP for 30 minutes."
		case 8:
			return "Lucky Charm activated. +10% crime success for 30 minutes."
		}
	}
	return ""
}
