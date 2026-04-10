package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	houseCost = 30
	maxPop    = 8
)

type SettlerHouse struct {
	loc       cbb.Coord
	pop       int
	woodStock int
	fishStock int
	needTimer time.Time
	taxTimer  time.Time
}

func NewSettlerHouse() cbb.Building { return &SettlerHouse{} }

func (h *SettlerHouse) GetLoc() cbb.Coord     { return h.loc }
func (h *SettlerHouse) Sprite() *ebiten.Image { return houseSprite }

func (h *SettlerHouse) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) != Water && aw.gold >= houseCost
}

func (h *SettlerHouse) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	h.needTimer = time.Now().Add(2 * time.Minute) // grace period before first need check
	h.taxTimer = time.Now().Add(30 * time.Second)
	aw := world.(*annoWorld)
	aw.gold -= houseCost
	aw.houses = append(aw.houses, h)
	return nil
}

func (h *SettlerHouse) CanUpdate(t time.Time) bool {
	return t.After(h.needTimer) || t.After(h.taxTimer)
}

func (h *SettlerHouse) Update(world cbb.World) {
	aw := world.(*annoWorld)
	now := time.Now()

	if now.After(h.needTimer) {
		h.needTimer = now.Add(60 * time.Second)
		if h.woodStock > 0 && h.fishStock > 0 {
			h.woodStock--
			h.fishStock--
			if h.pop < maxPop {
				h.pop++
				aw.population++
			}
		} else if h.pop > 0 {
			h.pop--
			aw.population--
		}
	}

	if now.After(h.taxTimer) {
		h.taxTimer = now.Add(30 * time.Second)
		aw.gold += h.pop / 2
	}
}

// needsDelivery reports whether the house is out of any good.
func (h *SettlerHouse) needsDelivery() bool {
	return h.woodStock == 0 || h.fishStock == 0
}

// deliver adds goods to the house's local stock.
func (h *SettlerHouse) deliver(wood, fish int) {
	h.woodStock += wood
	h.fishStock += fish
}
