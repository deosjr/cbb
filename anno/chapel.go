package main

import (
	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const chapelCost = 100

type Chapel struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int
}

func NewChapel() cbb.Building { return &Chapel{} }

func (c *Chapel) GetLoc() cbb.Coord                        { return c.loc }
func (c *Chapel) Sprite() *ebiten.Image                     { return chapelSprite }
func (c *Chapel) AccessPoint() cbb.Coord                    { return c.accessPt }
func (c *Chapel) GetFootprintSprite() (*ebiten.Image, int)  { return c.isoSprite, c.isoFootH }

func (c *Chapel) SetRotation(r int) {
	c.rotation = r
	c.isoSprite, c.isoFootH = cbb.NewIsoBoxSpriteMultiRotated(chapelWall, chapelRoof, chapelWallH, 2, 1, r)
	sw, sh := 2, 1
	if r%2 == 1 {
		sw, sh = 1, 2
	}
	c.accessPt = cbb.BuildingAccessPoint(c.loc, sw, sh, r)
}

func (c *Chapel) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) != Water && aw.gold >= chapelCost
}

func (c *Chapel) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	c.loc = loc
	c.SetRotation(0)
	aw := world.(*annoWorld)
	aw.gold -= chapelCost
	return nil
}
