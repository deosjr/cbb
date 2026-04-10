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
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) != Water
}

func (w *Warehouse) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	w.loc = loc
	aw := world.(*annoWorld)
	aw.warehouseLoc = loc
	return []cbb.Unit{&MarketCarrier{
		loc:  loc,
		home: loc,
		ts:   time.Now(),
	}}
}

// MarketCarrier walks from the warehouse to houses that need goods and back.
// It pathfinds on the road network, so houses must be connected by roads.
type MarketCarrier struct {
	loc       cbb.Coord
	home      cbb.Coord
	target    *SettlerHouse
	route     []cbb.Coord
	carryWood int
	carryFish int
	ts        time.Time
}

func (c *MarketCarrier) GetLoc() cbb.Coord         { return c.loc }
func (c *MarketCarrier) Sprite() *ebiten.Image     { return carrierSprite }
func (c *MarketCarrier) CanUpdate(t time.Time) bool { return t.After(c.ts) }

func (c *MarketCarrier) Update(world cbb.World) {
	c.ts = time.Now().Add(time.Second / 2)
	aw := world.(*annoWorld)

	// Advance one step if moving
	if len(c.route) > 1 {
		c.loc = c.route[len(c.route)-2]
		c.route = c.route[:len(c.route)-1]
		return
	}

	// Arrived somewhere (or start)
	if c.target != nil {
		// Deliver to house
		c.target.deliver(c.carryWood, c.carryFish)
		c.carryWood, c.carryFish = 0, 0
		c.target = nil
		// Return to warehouse
		if route, err := cbb.FindRoute(world.Roads(), c.loc, c.home); err == nil {
			c.route = route
		}
		return
	}

	// At warehouse: load up
	c.carryWood, c.carryFish = 0, 0
	if aw.warehouse.Take(Wood, 2) {
		c.carryWood = 2
	}
	if aw.warehouse.Take(Fish, 2) {
		c.carryFish = 2
	}
	if c.carryWood == 0 && c.carryFish == 0 {
		return // nothing to deliver, wait
	}

	// Find nearest house that needs goods and is reachable by road
	for _, house := range aw.houses {
		if !house.needsDelivery() {
			continue
		}
		if route, err := cbb.FindRoute(world.Roads(), c.loc, house.loc); err == nil {
			c.route = route
			c.target = house
			return
		}
	}

	// No reachable needy house: return goods
	aw.warehouse.Add(Wood, c.carryWood)
	aw.warehouse.Add(Fish, c.carryFish)
	c.carryWood, c.carryFish = 0, 0
}
