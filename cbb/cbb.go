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

// Task represents a unit of work: deliver a resource to a building.
type Task interface {
	Destination() Building
	Resource() Resource
	Amount() int
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

// Receiver is implemented by buildings that accept deliveries from units.
// When a unit arrives at its destination, Receive is called if the building implements this.
type Receiver interface {
	Receive(Task, World)
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
	AddTask(Task)
	ClaimTask() (Task, bool)
	Roads() *TileMap
	Tilemap() *TileMap
}

// Placeable is optionally implemented by buildings to restrict valid placement locations.
// CanPlace is called before WhenPlaced; returning false blocks placement.
type Placeable interface {
	CanPlace(Coord, World) bool
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
	Key     ebiten.Key
	Sprite  *ebiten.Image // cursor sprite; road options may omit this
	NewFunc func() Building
}
