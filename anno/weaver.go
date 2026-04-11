package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const weaverCost = 50

// weavingDur is the processing time per batch (2 Wool → 1 Cloth).
// Anno: 1.95t/min → ~31s per unit, compressed by gameSpeed.
const weavingDur = 31 * time.Second / gameSpeed

type WeavingHut struct {
	loc        cbb.Coord
	rotation   int
	accessPt   cbb.Coord
	isoSprite  *ebiten.Image
	isoFootH   int
	inputs     *cbb.Inventory // Wool waiting to be processed
	stockpile  *cbb.Inventory // finished Cloth, collected by WarehouseCart
	processing bool
	doneAt     time.Time
	ts         time.Time
}

func NewWeavingHut() cbb.Building {
	return &WeavingHut{
		inputs:    cbb.NewInventory(),
		stockpile: cbb.NewInventory(),
	}
}

func (h *WeavingHut) GetLoc() cbb.Coord                        { return h.loc }
func (h *WeavingHut) Sprite() *ebiten.Image                     { return weaverSprite }
func (h *WeavingHut) Stockpile() *cbb.Inventory                 { return h.stockpile }
func (h *WeavingHut) AccessPoint() cbb.Coord                    { return h.accessPt }
func (h *WeavingHut) GetFootprintSprite() (*ebiten.Image, int)  { return h.isoSprite, h.isoFootH }

func (h *WeavingHut) SetRotation(r int) {
	h.rotation = r
	h.isoSprite, h.isoFootH = cbb.NewIsoBoxSpriteMultiRotated(weaverWall, weaverRoof, weaverWallH, 2, 2, r)
	h.accessPt = cbb.BuildingAccessPoint(h.loc, 2, 2, r)
}

func (h *WeavingHut) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) != Water && aw.gold >= weaverCost
}

func (h *WeavingHut) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	h.SetRotation(0)
	h.ts = time.Now()
	aw := world.(*annoWorld)
	aw.gold -= weaverCost
	aw.producers = append(aw.producers, h)
	return []cbb.Unit{&ProcessorCart{
		loc:          loc,
		home:         loc,
		homeBuilding: h,
		needs:        map[Good]int{Wool: 2},
		inputs:       h.inputs,
		ts:           time.Now(),
		sprite:       processorCartSprite,
	}}
}

// WeavingHut implements cbb.Updatable: it checks for inputs and runs the
// production timer independently of its ProcessorCart.
func (h *WeavingHut) CanUpdate(t time.Time) bool { return t.After(h.ts) }

func (h *WeavingHut) Update(_ cbb.World) {
	h.ts = time.Now().Add(time.Second)
	if h.processing {
		if time.Now().After(h.doneAt) {
			h.stockpile.Add(Cloth, 1)
			h.processing = false
		}
		return
	}
	if h.inputs.Take(Wool, 2) {
		h.processing = true
		h.doneAt = time.Now().Add(weavingDur)
	}
}

// -------------------------------------------------------------------
// ProcessorCart: fetches inputs from the warehouse for a processor building.
// -------------------------------------------------------------------

type cartState int

const (
	cartWaiting     cartState = iota // at home, waiting until inputs run low
	cartToWarehouse                  // en route to warehouse to collect inputs
	cartToHome                       // carrying goods, returning to processor
)

type ProcessorCart struct {
	loc          cbb.Coord
	home         cbb.Coord       // fallback home if homeBuilding is nil
	homeBuilding cbb.Accessible  // optional; overrides home for return routing
	route        []cbb.Coord
	needs        map[Good]int   // e.g. {Wool: 2}
	inputs       *cbb.Inventory // processor's input stockpile (shared pointer)
	carrying     map[Good]int
	state        cartState
	ts           time.Time
	sprite       *ebiten.Image
}

func (c *ProcessorCart) homeCoord() cbb.Coord {
	if c.homeBuilding != nil {
		return c.homeBuilding.AccessPoint()
	}
	return c.home
}

func (c *ProcessorCart) GetLoc() cbb.Coord         { return c.loc }
func (c *ProcessorCart) Sprite() *ebiten.Image      { return c.sprite }
func (c *ProcessorCart) CanUpdate(t time.Time) bool { return t.After(c.ts) }

func (c *ProcessorCart) Update(world cbb.World) {
	c.ts = time.Now().Add(500 * time.Millisecond)
	aw := world.(*annoWorld)

	// Advance one step along the current route.
	if len(c.route) > 1 {
		c.loc = c.route[len(c.route)-2]
		c.route = c.route[:len(c.route)-1]
		return
	}

	switch c.state {
	case cartWaiting:
		// Check whether the processor is short on any input.
		needsRestock := false
		for g, n := range c.needs {
			if !c.inputs.Has(g, n) {
				needsRestock = true
				_ = g
				break
			}
		}
		if !needsRestock {
			c.ts = time.Now().Add(2 * time.Second)
			return
		}
		if !aw.warehousePlaced {
			c.ts = time.Now().Add(5 * time.Second)
			return
		}
		warehouseDest := aw.warehouseLoc
		if aw.warehouseBuilding != nil {
			warehouseDest = aw.warehouseBuilding.AccessPoint()
		}
		route, err := cbb.FindRoute(world.Roads(), c.loc, warehouseDest)
		if err != nil {
			c.ts = time.Now().Add(5 * time.Second)
			return
		}
		c.route = route
		c.state = cartToWarehouse

	case cartToWarehouse:
		// Arrived at warehouse — take needed goods.
		if c.carrying == nil {
			c.carrying = make(map[Good]int)
		}
		anyTaken := false
		for g, n := range c.needs {
			if !c.inputs.Has(g, n) {
				if aw.warehouse.Take(g, n) {
					c.carrying[g] = n
					anyTaken = true
				}
			}
		}
		if !anyTaken {
			c.state = cartWaiting
			c.ts = time.Now().Add(5 * time.Second)
			return
		}
		route, err := cbb.FindRoute(world.Roads(), c.loc, c.homeCoord())
		if err != nil {
			// Can't get home: return goods to warehouse.
			for g, n := range c.carrying {
				aw.warehouse.Add(g, n)
			}
			c.carrying = nil
			c.state = cartWaiting
			return
		}
		c.route = route
		c.state = cartToHome

	case cartToHome:
		// Arrived at processor — deposit into its input stockpile.
		for g, n := range c.carrying {
			c.inputs.Add(g, n)
		}
		c.carrying = nil
		c.state = cartWaiting
	}
}
