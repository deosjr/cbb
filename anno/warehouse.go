package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

type Warehouse struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int
}

func NewWarehouse() cbb.Building { return &Warehouse{} }

func (w *Warehouse) GetLoc() cbb.Coord                        { return w.loc }
func (w *Warehouse) Sprite() *ebiten.Image                     { return warehouseSprite }
func (w *Warehouse) AccessPoint() cbb.Coord                    { return w.accessPt }
func (w *Warehouse) GetFootprintSprite() (*ebiten.Image, int)  { return w.isoSprite, w.isoFootH }

func (w *Warehouse) SetRotation(r int) {
	w.rotation = r
	w.isoSprite, w.isoFootH = cbb.NewIsoBoxSpriteMultiRotated(warehouseWall, warehouseRoof, warehouseWallH, 2, 3, r)
	sw, sh := 2, 3
	if r%2 == 1 {
		sw, sh = 3, 2
	}
	w.accessPt = cbb.BuildingAccessPoint(w.loc, sw, sh, r)
}

func (w *Warehouse) CanPlace(loc cbb.Coord, world cbb.World) bool {
	return world.(*annoWorld).terrainAt(loc) != Water
}

func (w *Warehouse) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	w.loc = loc
	w.SetRotation(0)
	aw := world.(*annoWorld)
	aw.warehouseLoc = loc
	aw.warehouseBuilding = w
	aw.warehousePlaced = true
	return []cbb.Unit{&WarehouseCart{
		loc:          loc,
		home:         loc,
		homeBuilding: w,
		state:        wcIdle,
		ts:           time.Now(),
		sprite:       carrierSprite,
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
	loc          cbb.Coord
	home         cbb.Coord       // fallback home if homeBuilding is nil
	homeBuilding cbb.Accessible  // optional; overrides home for return routing
	route        []cbb.Coord
	carrying     map[Good]int
	target       Producer
	state        wcState
	ts           time.Time
	sprite       *ebiten.Image
}

func (c *WarehouseCart) homeCoord() cbb.Coord {
	if c.homeBuilding != nil {
		return c.homeBuilding.AccessPoint()
	}
	return c.home
}

func producerDest(p Producer) cbb.Coord {
	if a, ok := p.(cbb.Accessible); ok {
		return a.AccessPoint()
	}
	return p.GetLoc()
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
		// Route from the warehouse access point, not c.loc. The interior tile is
		// impassable in Roads (building footprint), so we always depart from the door.
		// This also handles the initial spawn where c.loc is still the interior tile.
		home := c.homeCoord()
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
			route, err := cbb.FindRoute(world.Roads(), home, producerDest(p))
			if err != nil {
				continue
			}
			c.loc = home // step to the door before departing
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
		route, err := cbb.FindRoute(world.Roads(), c.loc, c.homeCoord())
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
