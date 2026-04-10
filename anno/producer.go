package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

// workerState drives the out-and-back cycle shared by all extractor workers.
type workerState int

const (
	stateFinding    workerState = iota
	stateGoingOut
	stateHarvesting // waiting at the resource site
	stateReturning
)

// ExtractorWorker is the shared unit for Fisher, Hunter, Forester, and Sheep Farm.
// Each tick it either walks toward a resource tile, waits while harvesting,
// or walks home and deposits one unit into the building's local stockpile.
type ExtractorWorker struct {
	loc        cbb.Coord
	home       cbb.Coord
	stockpile  *cbb.Inventory // pointer to the spawning building's stockpile
	output     Good
	route      []cbb.Coord
	state      workerState
	harvestEnd time.Time
	harvestDur time.Duration
	ts         time.Time
	sprite     *ebiten.Image
	// findTarget returns the terrain tile to walk toward from the given origin.
	findTarget func(*annoWorld, cbb.Coord) (cbb.Coord, bool)
}

func (w *ExtractorWorker) GetLoc() cbb.Coord          { return w.loc }
func (w *ExtractorWorker) Sprite() *ebiten.Image       { return w.sprite }
func (w *ExtractorWorker) CanUpdate(t time.Time) bool  { return t.After(w.ts) }

func (w *ExtractorWorker) Update(world cbb.World) {
	w.ts = time.Now().Add(500 * time.Millisecond)
	aw := world.(*annoWorld)

	switch w.state {
	case stateFinding:
		target, ok := w.findTarget(aw, w.home)
		if !ok {
			w.ts = time.Now().Add(5 * time.Second)
			return
		}
		route, err := cbb.FindRoute(world.Tilemap(), w.loc, target)
		if err != nil {
			w.ts = time.Now().Add(5 * time.Second)
			return
		}
		w.route = route
		w.state = stateGoingOut

	case stateGoingOut:
		if len(w.route) <= 1 {
			w.harvestEnd = time.Now().Add(w.harvestDur)
			w.state = stateHarvesting
			return
		}
		w.loc = w.route[len(w.route)-2]
		w.route = w.route[:len(w.route)-1]

	case stateHarvesting:
		if time.Now().Before(w.harvestEnd) {
			return
		}
		route, err := cbb.FindRoute(world.Tilemap(), w.loc, w.home)
		if err != nil {
			w.state = stateFinding
			return
		}
		w.route = route
		w.state = stateReturning

	case stateReturning:
		if len(w.route) <= 1 {
			w.stockpile.Add(w.output, 1)
			w.state = stateFinding
			return
		}
		w.loc = w.route[len(w.route)-2]
		w.route = w.route[:len(w.route)-1]
	}
}
