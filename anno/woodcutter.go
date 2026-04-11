package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const foresterCost = 30

type ForesterHut struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	stockpile *cbb.Inventory
}

func NewForesterHut() cbb.Building { return &ForesterHut{stockpile: cbb.NewInventory()} }

func (h *ForesterHut) GetLoc() cbb.Coord        { return h.loc }
func (h *ForesterHut) Sprite() *ebiten.Image     { return foresterSprite }
func (h *ForesterHut) Stockpile() *cbb.Inventory { return h.stockpile }
func (h *ForesterHut) AccessPoint() cbb.Coord    { return h.accessPt }

func (h *ForesterHut) SetRotation(r int) {
	h.rotation = r
	h.accessPt = cbb.BuildingAccessPoint(h.loc, 2, 2, r)
}

func (h *ForesterHut) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) != Water && aw.gold >= foresterCost
}

func (h *ForesterHut) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	aw := world.(*annoWorld)
	aw.gold -= foresterCost
	aw.producers = append(aw.producers, h)
	// Anno: 2.8t/min → ~21s full cycle.
	return []cbb.Unit{&ExtractorWorker{
		home:         loc,
		homeBuilding: h,
		loc:          loc,
		stockpile:    h.stockpile,
		output:       Wood,
		harvestDur:   21*time.Second/gameSpeed - 2*time.Second,
		ts:           time.Now(),
		sprite:       foresterWorkerSprite,
		findTarget: func(aw *annoWorld, from cbb.Coord) (cbb.Coord, bool) {
			return findNearestTerrain(aw, from, Forest)
		},
	}}
}
