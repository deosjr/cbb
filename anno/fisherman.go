package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const fishermanCost = 50

type FishermanHut struct {
	loc cbb.Coord
}

func NewFishermanHut() cbb.Building { return &FishermanHut{} }

func (h *FishermanHut) GetLoc() cbb.Coord     { return h.loc }
func (h *FishermanHut) Sprite() *ebiten.Image { return fishermanSprite }

// CanPlace requires at least one adjacent water tile.
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
	return []cbb.Unit{&FishWorker{
		home: loc,
		loc:  loc,
		ts:   time.Now(),
	}}
}

type FishWorker struct {
	loc      cbb.Coord
	home     cbb.Coord
	route    []cbb.Coord
	state    workerState
	carrying int
	ts       time.Time
}

func (w *FishWorker) GetLoc() cbb.Coord         { return w.loc }
func (w *FishWorker) Sprite() *ebiten.Image     { return fishWorkerSprite }
func (w *FishWorker) CanUpdate(t time.Time) bool { return t.After(w.ts) }

func (w *FishWorker) Update(world cbb.World) {
	w.ts = time.Now().Add(time.Second / 2)
	aw := world.(*annoWorld)

	switch w.state {
	case stateFinding:
		target, ok := findNearestFishingSpot(aw, w.home)
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
			if w.carrying > 0 {
				aw.warehouse.Add(Fish, w.carrying)
				w.carrying = 0
			}
			w.state = stateFinding
			return
		}
		w.loc = w.route[len(w.route)-2]
		w.route = w.route[:len(w.route)-1]
	}
}
