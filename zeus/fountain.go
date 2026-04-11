package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const fountainPlaceCost = 60

// fountainWalkerSteps is how many tiles the fountain walker covers per patrol.
const fountainWalkerSteps = 30

// fountainWalkerCooldown is the rest period between fountain patrols.
const fountainWalkerCooldown = 15 * time.Second / gameSpeed

// Fountain spawns a hygiene walker that provides ServiceHygiene coverage.
// Houses need this service to evolve from tier 1 (House) to tier 2 (Townhouse).
type Fountain struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int
}

func NewFountain() cbb.Building { return &Fountain{} }

func (f *Fountain) GetLoc() cbb.Coord                       { return f.loc }
func (f *Fountain) Sprite() *ebiten.Image                   { return fountainSprite }
func (f *Fountain) AccessPoint() cbb.Coord                  { return f.accessPt }
func (f *Fountain) GetFootprintSprite() (*ebiten.Image, int) { return f.isoSprite, f.isoFootH }

func (f *Fountain) SetRotation(r int) {
	f.rotation = r
	f.isoSprite, f.isoFootH = cbb.NewIsoBoxSpriteMultiRotated(fountainWall, fountainRoof, fountainWallH, 1, 1, r)
	f.accessPt = cbb.BuildingAccessPoint(f.loc, 1, 1, r)
}

func (f *Fountain) CanPlace(loc cbb.Coord, world cbb.World) bool {
	zw := world.(*zeusWorld)
	return zw.terrainAt(loc) != Water && zw.gold >= fountainPlaceCost
}

func (f *Fountain) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	f.loc = loc
	f.SetRotation(0)
	zw := world.(*zeusWorld)
	zw.gold -= fountainPlaceCost
	walker := NewServiceWalker(
		f.accessPt,
		ServiceHygiene,
		fountainWalkerSprite,
		fountainWalkerSteps,
		fountainWalkerCooldown,
	)
	return []cbb.Unit{walker}
}
