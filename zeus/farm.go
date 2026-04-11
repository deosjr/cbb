package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const farmPlaceCost = 40

// farmProductionInterval is how often a Wheat Farm produces one unit of wheat.
const farmProductionInterval = 30 * time.Second / gameSpeed

// WheatFarm is a timer-based food producer. Every farmProductionInterval it
// deposits one unit of Wheat directly into the world's granary. No worker unit
// is spawned — the farm represents the whole field + worker household.
//
// In a fuller implementation this would spawn a worker carrying grain to the
// granary, but for the early-game demo the direct-deposit model keeps things
// simple and focuses attention on the service-walker mechanics.
type WheatFarm struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int
	ts        time.Time
}

func NewWheatFarm() cbb.Building { return &WheatFarm{} }

func (f *WheatFarm) GetLoc() cbb.Coord                       { return f.loc }
func (f *WheatFarm) Sprite() *ebiten.Image                   { return farmSprite }
func (f *WheatFarm) AccessPoint() cbb.Coord                  { return f.accessPt }
func (f *WheatFarm) GetFootprintSprite() (*ebiten.Image, int) { return f.isoSprite, f.isoFootH }

func (f *WheatFarm) CanUpdate(t time.Time) bool { return t.After(f.ts) }

func (f *WheatFarm) Update(world cbb.World) {
	f.ts = time.Now().Add(farmProductionInterval)
	zw := world.(*zeusWorld)
	zw.granary.Add(Wheat, 1)
}

func (f *WheatFarm) SetRotation(r int) {
	f.rotation = r
	f.isoSprite, f.isoFootH = cbb.NewIsoBoxSpriteMultiRotated(farmWall, farmRoof, farmWallH, 2, 2, r)
	f.accessPt = cbb.BuildingAccessPoint(f.loc, 2, 2, r)
}

func (f *WheatFarm) CanPlace(loc cbb.Coord, world cbb.World) bool {
	zw := world.(*zeusWorld)
	return zw.terrainAt(loc) == Plains && zw.gold >= farmPlaceCost
}

func (f *WheatFarm) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	f.loc = loc
	f.SetRotation(0)
	f.ts = time.Now().Add(farmProductionInterval)
	zw := world.(*zeusWorld)
	zw.gold -= farmPlaceCost
	return nil // farm itself is Updatable; no separate unit needed
}
