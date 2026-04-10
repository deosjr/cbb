package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const fishermanCost = 50

type FishermanHut struct {
	loc       cbb.Coord
	stockpile *cbb.Inventory
}

func NewFishermanHut() cbb.Building { return &FishermanHut{stockpile: cbb.NewInventory()} }

func (h *FishermanHut) GetLoc() cbb.Coord        { return h.loc }
func (h *FishermanHut) Sprite() *ebiten.Image     { return fishermanSprite }
func (h *FishermanHut) Stockpile() *cbb.Inventory { return h.stockpile }

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
	aw := world.(*annoWorld)
	aw.gold -= fishermanCost
	aw.producers = append(aw.producers, h)
	// Anno: 0.75t/min → ~80s full cycle. harvestDur = 80s/gameSpeed minus avg walk.
	return []cbb.Unit{&ExtractorWorker{
		home:       loc,
		loc:        loc,
		stockpile:  h.stockpile,
		output:     Food,
		harvestDur: 80*time.Second/gameSpeed - 4*time.Second,
		ts:         time.Now(),
		sprite:     fishWorkerSprite,
		findTarget: func(aw *annoWorld, from cbb.Coord) (cbb.Coord, bool) {
			return findNearestFishingSpot(aw, from)
		},
	}}
}
