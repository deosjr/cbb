package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const hunterPlaceCost = 40

// hunterProductionInterval is how often a Hunting Lodge produces one unit of meat.
// Slightly slower than farms to reflect the hunting cycle.
const hunterProductionInterval = 40 * time.Second / gameSpeed

// HuntingLodge is a timer-based meat producer that must be placed adjacent to
// or near Forest tiles to represent viable hunting grounds.
// Like WheatFarm it deposits directly into the world granary for simplicity.
type HuntingLodge struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int
	ts        time.Time
}

func NewHuntingLodge() cbb.Building { return &HuntingLodge{} }

func (h *HuntingLodge) GetLoc() cbb.Coord                       { return h.loc }
func (h *HuntingLodge) Sprite() *ebiten.Image                   { return hunterSprite }
func (h *HuntingLodge) AccessPoint() cbb.Coord                  { return h.accessPt }
func (h *HuntingLodge) GetFootprintSprite() (*ebiten.Image, int) { return h.isoSprite, h.isoFootH }

func (h *HuntingLodge) CanUpdate(t time.Time) bool { return t.After(h.ts) }

func (h *HuntingLodge) Update(world cbb.World) {
	h.ts = time.Now().Add(hunterProductionInterval)
	zw := world.(*zeusWorld)
	zw.granary.Add(Meat, 1)
}

func (h *HuntingLodge) SetRotation(r int) {
	h.rotation = r
	h.isoSprite, h.isoFootH = cbb.NewIsoBoxSpriteMultiRotated(hunterWall, hunterRoof, hunterWallH, 1, 1, r)
	h.accessPt = cbb.BuildingAccessPoint(h.loc, 1, 1, r)
}

func (h *HuntingLodge) CanPlace(loc cbb.Coord, world cbb.World) bool {
	zw := world.(*zeusWorld)
	if zw.gold < hunterPlaceCost {
		return false
	}
	// Must be near forest: at least one adjacent tile is Forest.
	for _, adj := range adjacentCoords(loc) {
		if zw.terrainAt(adj) == Forest {
			return true
		}
	}
	return false
}

func (h *HuntingLodge) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	h.SetRotation(0)
	h.ts = time.Now().Add(hunterProductionInterval)
	zw := world.(*zeusWorld)
	zw.gold -= hunterPlaceCost
	return nil
}
