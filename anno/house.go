package main

import (
	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	houseCost     = 30
	maxPioneerPop = 2
)

// PioneerHouse is a pure data struct. Population and tax are managed centrally
// by PopulationTick rather than per-house timers.
type PioneerHouse struct {
	loc cbb.Coord
	pop int
}

func NewPioneerHouse() cbb.Building { return &PioneerHouse{} }

func (h *PioneerHouse) GetLoc() cbb.Coord     { return h.loc }
func (h *PioneerHouse) Sprite() *ebiten.Image { return houseSprite }

func (h *PioneerHouse) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) != Water && aw.gold >= houseCost
}

func (h *PioneerHouse) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	aw := world.(*annoWorld)
	aw.gold -= houseCost
	aw.houses = append(aw.houses, h)
	return nil
}
