package cbb

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Resource is implemented by game-defined good types, typically an enum.
//
//	type Good int
//	const (Wood Good = iota; Stone; Grain)
//	func (g Good) Name() string { ... }
type Resource interface {
	Name() string
}

// Located is any entity that occupies a tile.
type Located interface {
	GetLoc() Coord
}

// Updatable is implemented by buildings and units that have periodic logic.
// Buildings opt in by implementing this; it is not required.
type Updatable interface {
	CanUpdate(time.Time) bool
	Update(World)
}

// Building is placed on the map and optionally spawns units when placed.
type Building interface {
	Located
	WhenPlaced(Coord, World) []Unit
	Sprite() *ebiten.Image
}

// Unit moves around the map and carries out tasks.
type Unit interface {
	Located
	Updatable
	Sprite() *ebiten.Image
}

// World is passed to Update and WhenPlaced, providing access to shared game state.
// Games extend this by embedding BaseWorld in their own world struct.
type World interface {
	Roads() *TileMap
	Tilemap() *TileMap
}

// Placeable is optionally implemented by buildings to restrict valid placement locations.
// CanPlace is called before WhenPlaced; returning false blocks placement.
type Placeable interface {
	CanPlace(Coord, World) bool
}

// FootprintSpriteGetter is optionally implemented by buildings that provide a
// single combined sprite for their full footprint in isometric mode.
// The engine draws it once at the anchor tile instead of one sprite per tile.
// Returns the sprite and footH (effective footprint height after rotation),
// which is used to compute the left-shift when positioning the sprite.
type FootprintSpriteGetter interface {
	GetFootprintSprite() (*ebiten.Image, int)
}

// Rotatable is optionally implemented by buildings that support rotation.
// SetRotation is called by the engine after WhenPlaced with the current rotation
// (0=south, 1=west, 2=north, 3=east).
type Rotatable interface {
	SetRotation(int)
}

// Accessible is optionally implemented by buildings that have a single entry/exit
// tile. Workers use AccessPoint as their home coordinate for pathfinding.
type Accessible interface {
	AccessPoint() Coord
}

// BuildingAccessPoint computes the access tile for a building footprint.
// anchor is the top-left tile, w and h are the effective dimensions after
// applying rotation. rotation: 0=south, 1=west, 2=north, 3=east.
func BuildingAccessPoint(anchor Coord, w, h, rotation int) Coord {
	switch rotation % 4 {
	case 0: // south: tile below the bottom edge
		return Coord{anchor.X + float64(w/2), anchor.Y + float64(h)}
	case 1: // west: tile left of the left edge
		return Coord{anchor.X - 1, anchor.Y + float64(h/2)}
	case 2: // north: tile above the top edge
		return Coord{anchor.X + float64(w/2), anchor.Y - 1}
	case 3: // east: tile right of the right edge
		return Coord{anchor.X + float64(w), anchor.Y + float64(h/2)}
	}
	return anchor
}

// SelectionKind distinguishes tool types in the build menu.
type SelectionKind int

const (
	KindRoad     SelectionKind = iota
	KindBuilding SelectionKind = iota
)

// Option describes a player-selectable build tool.
type Option struct {
	Name    string
	Kind    SelectionKind
	Radius  int
	SizeW   int // footprint width in tiles; 0 treated as 1
	SizeH   int // footprint height in tiles; 0 treated as 1
	Key     ebiten.Key
	Sprite  *ebiten.Image // cursor sprite; road options may omit this
	NewFunc func() Building
}
