package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

// buildingCoord returns the access point of a building if it implements
// cbb.Accessible, otherwise falls back to GetLoc. Units should always route
// to/from access points so that pathfinding stays on the road network rather
// than trying to enter impassable building-interior tiles.
func buildingCoord(b cbb.Building) cbb.Coord {
	if a, ok := b.(cbb.Accessible); ok {
		return a.AccessPoint()
	}
	return b.GetLoc()
}

// Gatherer walks to a task destination, triggers delivery, then returns home.
type Gatherer struct {
	task Task
	home cbb.Building
	loc  cbb.Coord
	ts   time.Time
}

func (g *Gatherer) GetLoc() cbb.Coord     { return g.loc }
func (g *Gatherer) Sprite() *ebiten.Image { return unitSprite }

func (g *Gatherer) CanUpdate(t time.Time) bool { return t.After(g.ts) }

func (g *Gatherer) Update(w cbb.World) {
	g.ts = time.Now().Add(time.Second / 2)
	tw := w.(*taskWorld)

	if g.task == nil {
		if t, ok := tw.ClaimTask(); ok {
			g.task = t
		}
	}
	if g.task == nil {
		return
	}

	dest := buildingCoord(g.task.Destination())
	route, err := cbb.FindRoute(w.Roads(), g.loc, dest)
	if err != nil {
		tw.AddTask(g.task)
		g.task = nil
		return
	}

	// Advance one step along the route.
	if len(route) == 1 {
		g.loc = route[0]
	} else {
		g.loc = route[len(route)-2]
	}

	if g.loc != dest {
		return
	}

	// Arrived at destination.
	if r, ok := g.task.Destination().(Receiver); ok {
		r.Receive(g.task, tw)
	}
	homeDest := buildingCoord(g.home)
	if g.loc == homeDest {
		// Return trip complete.
		g.task = nil
	} else {
		// Head home.
		g.task = &task{dest: g.home}
	}
}
