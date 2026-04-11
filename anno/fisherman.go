package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const fishermanCost = 50

type FishermanHut struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int
	stockpile *cbb.Inventory
}

func NewFishermanHut() cbb.Building { return &FishermanHut{stockpile: cbb.NewInventory()} }

func (h *FishermanHut) GetLoc() cbb.Coord                         { return h.loc }
func (h *FishermanHut) Sprite() *ebiten.Image                      { return fishermanSprite }
func (h *FishermanHut) Stockpile() *cbb.Inventory                  { return h.stockpile }
func (h *FishermanHut) AccessPoint() cbb.Coord                     { return h.accessPt }
func (h *FishermanHut) GetFootprintSprite() (*ebiten.Image, int)   { return h.isoSprite, h.isoFootH }

func (h *FishermanHut) SetRotation(r int) {
	h.rotation = r
	h.isoSprite, h.isoFootH = cbb.NewIsoBoxSpriteMultiRotated(fishermanWall, fishermanRoof, fishermanWallH, 1, 1, r)
	h.accessPt = cbb.BuildingAccessPoint(h.loc, 1, 1, r)
}

func (h *FishermanHut) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	if aw.gold < fishermanCost {
		return false
	}
	for _, adj := range adjacentCoords(loc) {
		if aw.terrainAt(adj) == Water {
			return true
		}
	}
	return false
}

func (h *FishermanHut) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	h.SetRotation(0)
	aw := world.(*annoWorld)
	aw.gold -= fishermanCost
	aw.producers = append(aw.producers, h)
	// Anno: 0.75t/min → ~80s full cycle. harvestDur = 80s/gameSpeed minus avg walk.
	return []cbb.Unit{&ExtractorWorker{
		home:         loc,
		homeBuilding: h,
		loc:          loc,
		stockpile:    h.stockpile,
		output:       Food,
		harvestDur:   80*time.Second/gameSpeed - 4*time.Second,
		ts:           time.Now(),
		sprite:       fishWorkerSprite,
		findTarget: func(aw *annoWorld, from cbb.Coord) (cbb.Coord, bool) {
			return findNearestFishingSpot(aw, from)
		},
	}}
}
