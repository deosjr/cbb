package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

// task is the game's concrete implementation of Task.
type task struct {
	dest     cbb.Building
	resource cbb.Resource
	amount   int
}

func (t *task) Destination() cbb.Building { return t.dest }
func (t *task) Resource() cbb.Resource    { return t.resource }
func (t *task) Amount() int               { return t.amount }

// Producer generates tasks periodically and waits for a gatherer to collect them.
type Producer struct {
	loc      cbb.Coord
	accessPt cbb.Coord
	loaded   bool
	ts       time.Time
}

func NewProducer() cbb.Building { return &Producer{} }

func (p *Producer) GetLoc() cbb.Coord      { return p.loc }
func (p *Producer) AccessPoint() cbb.Coord { return p.accessPt }
func (p *Producer) Sprite() *ebiten.Image  { return producerSprite }

func (p *Producer) WhenPlaced(loc cbb.Coord, _ cbb.World) []cbb.Unit {
	p.loc = loc
	p.accessPt = cbb.BuildingAccessPoint(loc, 1, 1, 0)
	p.ts = time.Now()
	return nil
}

func (p *Producer) CanUpdate(t time.Time) bool {
	return !p.loaded && t.After(p.ts)
}

func (p *Producer) Update(w cbb.World) {
	p.loaded = true
	w.(*taskWorld).AddTask(&task{dest: p})
	p.ts = time.Now().Add(30 * time.Second)
}

// Receive is called by a gatherer on arrival. Clears the loaded flag so the
// producer can generate its next task after the cooldown.
func (p *Producer) Receive(_ Task, _ *taskWorld) {
	p.loaded = false
}

// Consumer spawns a gatherer when placed. Has no periodic update logic.
type Consumer struct {
	loc      cbb.Coord
	accessPt cbb.Coord
}

func NewConsumer() cbb.Building { return &Consumer{} }

func (c *Consumer) GetLoc() cbb.Coord      { return c.loc }
func (c *Consumer) AccessPoint() cbb.Coord { return c.accessPt }
func (c *Consumer) Sprite() *ebiten.Image  { return consumerSprite }

func (c *Consumer) WhenPlaced(loc cbb.Coord, _ cbb.World) []cbb.Unit {
	c.loc = loc
	c.accessPt = cbb.BuildingAccessPoint(loc, 1, 1, 0)
	return []cbb.Unit{&Gatherer{
		ts:   time.Now(),
		home: c,
		loc:  c.accessPt, // start at the door, not inside the building
	}}
}
