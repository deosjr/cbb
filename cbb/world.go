package cbb

// BaseWorld provides the standard World implementation.
// Embed this in a game-specific world struct to extend it with additional state
// (inventory, gold, terrain data, etc.) while satisfying the World interface.
//
//	type MyWorld struct {
//	    *cbb.BaseWorld
//	    gold int
//	}
type BaseWorld struct {
	tilemap *TileMap
	roads   *TileMap
}

// NewBaseWorld creates a BaseWorld from a tile map.
// Pass the returned value (or a struct embedding it) to NewGame.
func NewBaseWorld(tilemap *TileMap) *BaseWorld {
	return &BaseWorld{
		tilemap: tilemap,
		roads:   &TileMap{Tiles: map[Coord]Tile{}},
	}
}

func (w *BaseWorld) Roads() *TileMap {
	return w.roads
}

func (w *BaseWorld) Tilemap() *TileMap {
	return w.tilemap
}
