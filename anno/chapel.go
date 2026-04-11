package main

import (
	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const chapelCost = 100

type Chapel struct {
	loc      cbb.Coord
	rotation int
	accessPt cbb.Coord
}

func NewChapel() cbb.Building { return &Chapel{} }

func (c *Chapel) GetLoc() cbb.Coord      { return c.loc }
func (c *Chapel) Sprite() *ebiten.Image  { return chapelSprite }
func (c *Chapel) AccessPoint() cbb.Coord { return c.accessPt }

func (c *Chapel) SetRotation(r int) {
	c.rotation = r
	c.accessPt = cbb.BuildingAccessPoint(c.loc, 2, 1, r)
}

func (c *Chapel) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) != Water && aw.gold >= chapelCost
}

func (c *Chapel) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	c.loc = loc
	aw := world.(*annoWorld)
	aw.gold -= chapelCost
	return nil
}
