package main

import (
	"time"

	"github.com/deosjr/Pathfinding/path"
)

type Updatable interface {
	CanUpdate(time.Time) bool
	Update()
}

type Building interface {
	WhenPlaced(coord)
}

type producer struct {
	timestamp time.Time
	loaded    bool
	loc       coord
}

func (p *producer) WhenPlaced(loc coord) {
	p.timestamp = time.Now()
	p.loc = loc
}

func (p *producer) CanUpdate(t time.Time) bool {
	if p.loaded {
		return false
	}
	return t.After(p.timestamp)
}

func (p *producer) Update() {
	p.loaded = true
	tasks = append(tasks, &task{home: p, location: p.loc})
	dt := 30 * time.Second
	p.timestamp = time.Now().Add(dt)
}

type consumer struct {
	unit *gatherer
	loc  coord
}

func (c *consumer) WhenPlaced(loc coord) {
	c.loc = loc
	c.unit = &gatherer{
		timestamp: time.Now(),
		home:      c,
		loc:       loc,
	}
}

type Unit interface {
	GetLoc() coord
}

type gatherer struct {
	timestamp time.Time
	task      Task
	home      Building
	loc       coord
}

func (g *gatherer) CanUpdate(t time.Time) bool {
	return t.After(g.timestamp)
}

func (g *gatherer) Update() {
	dt := time.Second / 2
	g.timestamp = time.Now().Add(dt)
	if g.task == nil && len(tasks) > 0 {
		g.task = tasks[0]
		tasks = tasks[1:]
	}
	if g.task == nil {
		return
	}
	route, err := path.FindRoute(roads, g.loc, g.task.GetLoc())
	if err != nil {
		tasks = append(tasks, g.task)
		g.task = nil
		return
	}
	var next coord
	if len(route) == 1 {
		next = route[0].(coord)
	} else {
		next = route[len(route)-2].(coord)
	}
	g.loc = next
	if g.loc == g.task.GetLoc() && g.loc != g.home.(*consumer).loc {
		// destination reached
		g.task.GetHome().(*producer).loaded = false
		g.task = &task{location: g.home.(*consumer).loc}
		return
	}
	if g.loc == g.home.(*consumer).loc && g.task.GetLoc() == g.home.(*consumer).loc {
		// return trip finished
		g.task = nil
	}
}

func (g *gatherer) GetLoc() coord {
	return g.loc
}

type Task interface {
	GetHome() Building
	GetLoc() coord
}

// get vs return tasks as types?
type task struct {
	home     Building
	location coord
}

func (t *task) GetLoc() coord {
	return t.location
}

func (t *task) GetHome() Building {
	return t.home
}
