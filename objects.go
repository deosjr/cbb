package main

import (
	"time"

	"github.com/deosjr/Pathfinding/path"
)

type Instantiate func() Building

type Updatable interface {
	CanUpdate(time.Time) bool
	Update()
}

type updatable struct {
	timestamp time.Time
}

func (u updatable) CanUpdate(t time.Time) bool {
	return t.After(u.timestamp)
}

type Building interface {
	WhenPlaced(loc coord) (upd []Updatable, dyn []Updatable)
	GetLoc() coord
}

type building struct {
	loc coord
}

func (b building) GetLoc() coord {
	return b.loc
}

type producer struct {
	updatable
	building
	loaded bool
}

func NewProducer() Building {
	return &producer{}
}

func (p *producer) WhenPlaced(loc coord) ([]Updatable, []Updatable) {
	p.timestamp = time.Now()
	p.loc = loc
	return []Updatable{p}, nil
}

func (p *producer) CanUpdate(t time.Time) bool {
	if p.loaded {
		return false
	}
	return t.After(p.timestamp)
}

func (p *producer) Update() {
	p.loaded = true
	tasks = append(tasks, &task{dest: p})
	// TODO: timer should start from moment loaded=false instead
	dt := 30 * time.Second
	p.timestamp = time.Now().Add(dt)
}

type consumer struct {
	building
	unit *gatherer
}

func NewConsumer() Building {
	return &consumer{}
}

func (c *consumer) WhenPlaced(loc coord) ([]Updatable, []Updatable) {
	c.loc = loc
	c.unit = &gatherer{
		updatable: updatable{
			timestamp: time.Now(),
		},
		home: c,
		loc:  loc,
	}
	return []Updatable{c.unit}, []Updatable{c.unit}
}

type Unit interface {
	GetLoc() coord
}

type gatherer struct {
	updatable
	task Task
	home Building
	loc  coord
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
	route, err := path.FindRoute(roads, g.loc, g.task.GetDestination().GetLoc())
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
	if g.loc == g.task.GetDestination().GetLoc() {
		if g.loc == g.home.GetLoc() {
			// return trip finished
			g.task = nil
			return
		}
		// destination reached
		g.task.GetDestination().(*producer).loaded = false
		g.task = &task{dest: g.home}
		return
	}
}

func (g *gatherer) GetLoc() coord {
	return g.loc
}

type Task interface {
	GetDestination() Building
}

type task struct {
	dest Building
}

func (t *task) GetDestination() Building {
	return t.dest
}
