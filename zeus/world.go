package main

import (
	"image/color"
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

// ServiceType identifies which service a walker provides.
type ServiceType int

const (
	ServiceFood          ServiceType = iota // market lady walker
	ServiceHygiene                          // fountain walker
	ServiceEntertainment                    // theater troupe walker
	ServiceReligion                         // sanctuary priest walker
)

// coverageExpiry is how long a service visit remains valid.
// After this duration without a walker passing by, the house loses that service.
const coverageExpiry = 2 * time.Minute / gameSpeed

type zeusWorld struct {
	*cbb.BaseWorld
	terrain    map[cbb.Coord]Terrain
	granary    *cbb.Inventory // global food store; set/shared by Granary buildings
	gold       int
	population int
	houses     []*GreekHouse

	// coverage records the last time each (tile, service) pair was visited by
	// a service walker. Houses check this to determine tier eligibility.
	coverage map[cbb.Coord]map[ServiceType]time.Time
}

func newZeusWorld(tilemap *cbb.TileMap, terrain map[cbb.Coord]Terrain) *zeusWorld {
	return &zeusWorld{
		BaseWorld: cbb.NewBaseWorld(tilemap),
		terrain:   terrain,
		granary:   cbb.NewInventory(),
		gold:      1500,
		coverage:  map[cbb.Coord]map[ServiceType]time.Time{},
	}
}

func (w *zeusWorld) terrainAt(c cbb.Coord) Terrain {
	return w.terrain[c]
}

// RecordCoverage marks a tile as visited by the given service right now.
func (w *zeusWorld) RecordCoverage(c cbb.Coord, s ServiceType) {
	if w.coverage[c] == nil {
		w.coverage[c] = map[ServiceType]time.Time{}
	}
	w.coverage[c][s] = time.Now()
}

// HasCoverage reports whether a tile has been visited by the given service recently.
func (w *zeusWorld) HasCoverage(c cbb.Coord, s ServiceType) bool {
	if m, ok := w.coverage[c]; ok {
		if t, ok := m[s]; ok {
			return time.Since(t) < coverageExpiry
		}
	}
	return false
}

// isoBoxMulti computes a multi-tile iso box sprite and its footH for GetFootprintSprite.
// baseW/baseH are the unrotated dimensions; rotation swaps them at odd values.
func isoBoxMulti(wall, roof color.Color, wallH, baseW, baseH, rotation int) (*ebiten.Image, int) {
	w, h := baseW, baseH
	if rotation%2 == 1 {
		w, h = h, w
	}
	return cbb.NewIsoBoxSpriteMulti(wall, roof, wallH, w, h), h
}
