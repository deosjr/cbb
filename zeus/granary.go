package main

import (
	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const granaryPlaceCost = 100

// Granary is the central food storage building. When placed it becomes the
// world's active granary inventory — all food producers deposit here, and
// PopulationTick consumes from here.
//
// Only one granary is assumed in the early game. Future work could support
// multiple granaries with cart delivery between them.
type Granary struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int
}

func NewGranary() cbb.Building { return &Granary{} }

func (g *Granary) GetLoc() cbb.Coord                       { return g.loc }
func (g *Granary) Sprite() *ebiten.Image                   { return granarySprite }
func (g *Granary) AccessPoint() cbb.Coord                  { return g.accessPt }
func (g *Granary) GetFootprintSprite() (*ebiten.Image, int) { return g.isoSprite, g.isoFootH }

func (g *Granary) SetRotation(r int) {
	g.rotation = r
	g.isoSprite, g.isoFootH = cbb.NewIsoBoxSpriteMultiRotated(granaryWall, granaryRoof, granaryWallH, 2, 2, r)
	g.accessPt = cbb.BuildingAccessPoint(g.loc, 2, 2, r)
}

func (g *Granary) CanPlace(loc cbb.Coord, world cbb.World) bool {
	zw := world.(*zeusWorld)
	return zw.terrainAt(loc) != Water && zw.gold >= granaryPlaceCost
}

func (g *Granary) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	g.loc = loc
	g.SetRotation(0)
	zw := world.(*zeusWorld)
	zw.gold -= granaryPlaceCost
	// The world's granary inventory is already initialised; we keep using it
	// so any food produced before placement is preserved.
	return nil
}
