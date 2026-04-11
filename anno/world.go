package main

import "github.com/deosjr/tiles/cbb"

// gameSpeed compresses all production and consumption timers so the game
// is observable in real time. A value of 10 means 10 in-game minutes pass
// per real second, preserving the supply/demand ratios from the Anno data.
const gameSpeed = 10

// Producer is implemented by any building that has a local output stockpile
// the WarehouseCart should collect from.
type Producer interface {
	cbb.Located
	Stockpile() *cbb.Inventory
}

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
