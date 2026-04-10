package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const sheepFarmCost = 50

type SheepFarm struct {
	loc       cbb.Coord
	stockpile *cbb.Inventory
}

func NewSheepFarm() cbb.Building { return &SheepFarm{stockpile: cbb.NewInventory()} }

func (h *SheepFarm) GetLoc() cbb.Coord        { return h.loc }
func (h *SheepFarm) Sprite() *ebiten.Image     { return sheepFarmSprite }
func (h *SheepFarm) Stockpile() *cbb.Inventory { return h.stockpile }

func (h *SheepFarm) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) == Plains && aw.gold >= sheepFarmCost
}

func (h *SheepFarm) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	aw := world.(*annoWorld)
	aw.gold -= sheepFarmCost
	aw.producers = append(aw.producers, h)
	// Anno: 1.95t/min → ~31s full cycle.
	return []cbb.Unit{&ExtractorWorker{
		home:       loc,
		loc:        loc,
		stockpile:  h.stockpile,
		output:     Wool,
		harvestDur: 31*time.Second/gameSpeed - 2*time.Second,
		ts:         time.Now(),
		sprite:     sheepWorkerSprite,
		findTarget: func(aw *annoWorld, from cbb.Coord) (cbb.Coord, bool) {
			return findNearestTerrain(aw, from, Plains)
		},
	}}
}
