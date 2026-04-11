package main

import (
	"github.com/deosjr/tiles/cbb"
)

type annoWorld struct {
	*cbb.BaseWorld
	terrain           map[cbb.Coord]Terrain
	warehouse         *cbb.Inventory
	warehouseLoc      cbb.Coord
	warehouseBuilding cbb.Accessible // set after warehouse is placed+rotated
	warehousePlaced   bool
	gold              int
	population        int
	houses            []*PioneerHouse
	producers         []Producer
}

func newAnnoWorld(tilemap *cbb.TileMap, terrain map[cbb.Coord]Terrain) *annoWorld {
	return &annoWorld{
		BaseWorld: cbb.NewBaseWorld(tilemap),
		terrain:   terrain,
		warehouse: cbb.NewInventory(),
		gold:      2000,
	}
}

func (w *annoWorld) terrainAt(c cbb.Coord) Terrain {
	return w.terrain[c]
}

