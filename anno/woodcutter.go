package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const woodcutterCost = 50

type WoodcutterHut struct {
	loc cbb.Coord
}

func NewWoodcutterHut() cbb.Building { return &WoodcutterHut{} }

func (h *WoodcutterHut) GetLoc() cbb.Coord     { return h.loc }
func (h *WoodcutterHut) Sprite() *ebiten.Image { return woodcutterSprite }

func (h *WoodcutterHut) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) != Water && aw.gold >= woodcutterCost
}

func (h *WoodcutterHut) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	aw := world.(*annoWorld)
	aw.gold -= woodcutterCost
	return []cbb.Unit{&WoodWorker{
		home: loc,
		loc:  loc,
		ts:   time.Now(),
	}}
}

// workerState drives the out-and-back cycle for production workers.
type workerState int

const (
	stateFinding   workerState = iota
	stateGoingOut  workerState = iota
	stateReturning workerState = iota
)

type WoodWorker struct {
	loc      cbb.Coord
	home     cbb.Coord
	route    []cbb.Coord
	state    workerState
	carrying int
	ts       time.Time
}

func (w *WoodWorker) GetLoc() cbb.Coord         { return w.loc }
func (w *WoodWorker) Sprite() *ebiten.Image     { return woodWorkerSprite }
func (w *WoodWorker) CanUpdate(t time.Time) bool { return t.After(w.ts) }

func (w *WoodWorker) Update(world cbb.World) {
	w.ts = time.Now().Add(time.Second / 2)
	aw := world.(*annoWorld)

	switch w.state {
	case stateFinding:
		target, ok := findNearestTerrain(aw, w.home, Forest)
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
			// Arrived at forest tile: pick up wood
			w.carrying = 1
			route, err := cbb.FindRoute(world.Tilemap(), w.loc, w.home)
			if err != nil {
				w.state = stateFinding
				return
			}
			w.route = route
			w.state = stateReturning
			return
		}
		w.loc = w.route[len(w.route)-2]
		w.route = w.route[:len(w.route)-1]

	case stateReturning:
		if len(w.route) <= 1 {
			// Arrived home: deposit
			if w.carrying > 0 {
				aw.warehouse.Add(Wood, w.carrying)
				w.carrying = 0
			}
			w.state = stateFinding
			return
		}
		w.loc = w.route[len(w.route)-2]
		w.route = w.route[:len(w.route)-1]
	}
}
