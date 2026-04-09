package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

// Gatherer walks to a task destination, triggers delivery, then returns home.
type Gatherer struct {
	task cbb.Task
	home cbb.Building
	loc  cbb.Coord
	ts   time.Time
}

func (g *Gatherer) GetLoc() cbb.Coord     { return g.loc }
func (g *Gatherer) Sprite() *ebiten.Image { return unitSprite }

func (g *Gatherer) CanUpdate(t time.Time) bool { return t.After(g.ts) }

func (g *Gatherer) Update(w cbb.World) {
	g.ts = time.Now().Add(time.Second / 2)

	if g.task == nil {
		if t, ok := w.ClaimTask(); ok {
			g.task = t
		}
	}
	if g.task == nil {
		return
	}

	route, err := cbb.FindRoute(w.Roads(), g.loc, g.task.Destination().GetLoc())
	if err != nil {
		w.AddTask(g.task)
		g.task = nil
		return
	}

	// Advance one step along the route (route is destination-first, so second-to-last is next step).
	if len(route) == 1 {
		g.loc = route[0]
	} else {
		g.loc = route[len(route)-2]
	}

	if g.loc != g.task.Destination().GetLoc() {
		return
	}

	// Arrived at destination.
	if r, ok := g.task.Destination().(cbb.Receiver); ok {
		r.Receive(g.task, w)
	}
	if g.loc == g.home.GetLoc() {
		// Return trip complete.
		g.task = nil
	} else {
		// Head home.
		g.task = &task{dest: g.home}
	}
}
