package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

type Warehouse struct {
	loc cbb.Coord
}

func NewWarehouse() cbb.Building { return &Warehouse{} }

func (w *Warehouse) GetLoc() cbb.Coord     { return w.loc }
func (w *Warehouse) Sprite() *ebiten.Image { return warehouseSprite }

func (w *Warehouse) CanPlace(loc cbb.Coord, world cbb.World) bool {
	return world.(*annoWorld).terrainAt(loc) != Water
}

func (w *Warehouse) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	w.loc = loc
	aw := world.(*annoWorld)
	aw.warehouseLoc = loc
	aw.warehousePlaced = true
	return []cbb.Unit{&WarehouseCart{
		loc:    loc,
		home:   loc,
		state:  wcIdle,
		ts:     time.Now(),
		sprite: carrierSprite,
	}}
}

// -------------------------------------------------------------------
// WarehouseCart: roams the road network collecting goods from all
// registered producers and returning them to the warehouse.
// -------------------------------------------------------------------

type wcState int

const (
	wcIdle       wcState = iota
	wcToProducer         // heading to a producer to collect
	wcToHome             // carrying goods back to warehouse
)

type WarehouseCart struct {
	loc      cbb.Coord
	home     cbb.Coord
	route    []cbb.Coord
	carrying map[Good]int
	target   Producer
	state    wcState
	ts       time.Time
	sprite   *ebiten.Image
}

func (c *WarehouseCart) GetLoc() cbb.Coord         { return c.loc }
func (c *WarehouseCart) Sprite() *ebiten.Image      { return c.sprite }
func (c *WarehouseCart) CanUpdate(t time.Time) bool { return t.After(c.ts) }

func (c *WarehouseCart) Update(world cbb.World) {
	c.ts = time.Now().Add(500 * time.Millisecond)
	aw := world.(*annoWorld)

	// Advance one step along the current route.
	if len(c.route) > 1 {
		c.loc = c.route[len(c.route)-2]
		c.route = c.route[:len(c.route)-1]
		return
	}

	switch c.state {
	case wcIdle:
		// Find the first producer with goods that is reachable by road.
		for _, p := range aw.producers {
			hasGoods := false
			for _, g := range AllGoods {
				if p.Stockpile().Count(g) > 0 {
					hasGoods = true
					break
				}
			}
			if !hasGoods {
				continue
			}
			route, err := cbb.FindRoute(world.Roads(), c.loc, p.GetLoc())
			if err != nil {
				continue
			}
			c.route = route
			c.target = p
			c.state = wcToProducer
			return
		}
		c.ts = time.Now().Add(2 * time.Second)

	case wcToProducer:
		// Arrived at producer: collect everything in its stockpile.
		if c.carrying == nil {
			c.carrying = make(map[Good]int)
		}
		for _, g := range AllGoods {
			n := c.target.Stockpile().Count(g)
			if n > 0 {
				c.target.Stockpile().Take(g, n)
				c.carrying[g] += n
			}
		}
		c.target = nil
		route, err := cbb.FindRoute(world.Roads(), c.loc, c.home)
		if err != nil {
			c.state = wcIdle
			return
		}
		c.route = route
		c.state = wcToHome

	case wcToHome:
		// Arrived at warehouse: deposit all carried goods.
		for g, n := range c.carrying {
			if n > 0 {
				aw.warehouse.Add(g, n)
				c.carrying[g] = 0
			}
		}
		c.state = wcIdle
	}
}
