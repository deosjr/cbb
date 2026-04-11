package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const hunterCost = 30

type HunterLodge struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	stockpile *cbb.Inventory
}

func NewHunterLodge() cbb.Building { return &HunterLodge{stockpile: cbb.NewInventory()} }

func (h *HunterLodge) GetLoc() cbb.Coord        { return h.loc }
func (h *HunterLodge) Sprite() *ebiten.Image     { return hunterSprite }
func (h *HunterLodge) Stockpile() *cbb.Inventory { return h.stockpile }
func (h *HunterLodge) AccessPoint() cbb.Coord    { return h.accessPt }

func (h *HunterLodge) SetRotation(r int) {
	h.rotation = r
	h.accessPt = cbb.BuildingAccessPoint(h.loc, 1, 1, r)
}

func (h *HunterLodge) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) != Water && aw.gold >= hunterCost
}

func (h *HunterLodge) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	aw := world.(*annoWorld)
	aw.gold -= hunterCost
	aw.producers = append(aw.producers, h)
	// Anno: 2t/min → ~30s full cycle.
	return []cbb.Unit{&ExtractorWorker{
		home:         loc,
		homeBuilding: h,
		loc:          loc,
		stockpile:    h.stockpile,
		output:       Food,
		harvestDur:   30*time.Second/gameSpeed - 2*time.Second,
		ts:           time.Now(),
		sprite:       hunterWorkerSprite,
		findTarget: func(aw *annoWorld, from cbb.Coord) (cbb.Coord, bool) {
			return findNearestTerrain(aw, from, Forest)
		},
	}}
}
