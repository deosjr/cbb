package main

import (
	"math/rand"
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const sheepFarmCost = 50

type SheepFarm struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int
	stockpile *cbb.Inventory
}

func NewSheepFarm() cbb.Building { return &SheepFarm{stockpile: cbb.NewInventory()} }

func (h *SheepFarm) GetLoc() cbb.Coord                        { return h.loc }
func (h *SheepFarm) Sprite() *ebiten.Image                    { return sheepFarmSprite }
func (h *SheepFarm) Stockpile() *cbb.Inventory                { return h.stockpile }
func (h *SheepFarm) AccessPoint() cbb.Coord                   { return h.accessPt }
func (h *SheepFarm) GetFootprintSprite() (*ebiten.Image, int) { return h.isoSprite, h.isoFootH }

func (h *SheepFarm) SetRotation(r int) {
	h.rotation = r
	h.isoSprite, h.isoFootH = isoBoxMulti(sheepFarmWall, sheepFarmRoof, sheepFarmWallH, 2, 2, r)
	h.accessPt = cbb.BuildingAccessPoint(h.loc, 2, 2, r)
}

func (h *SheepFarm) CanPlace(loc cbb.Coord, world cbb.World) bool {
	aw := world.(*annoWorld)
	return aw.terrainAt(loc) == Plains && aw.gold >= sheepFarmCost
}

func (h *SheepFarm) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	h.SetRotation(0)
	aw := world.(*annoWorld)
	aw.gold -= sheepFarmCost
	aw.producers = append(aw.producers, h)

	// Spawn 3 sheep that wander the surrounding plains. Each deposits 1 Wool on
	// returning home; together they produce ~1.95t/min as per Anno data.
	units := make([]cbb.Unit, 3)
	for i := range units {
		units[i] = &SheepUnit{
			loc:       loc,
			home:      loc,
			stockpile: h.stockpile,
			ts:        time.Now().Add(time.Duration(i) * 3 * time.Second), // stagger starts
			rng:       rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))),
		}
	}
	return units
}

// sheepState drives each sheep's wander-and-return cycle.
type sheepState int

const (
	sheepWandering sheepState = iota
	sheepReturning
)

// sheepMaxSteps is how many terrain tiles a sheep wanders before heading home.
// Anno rate 1.95t/min for 3 sheep ≈ one wool per sheep per ~31s/gameSpeed.
// At 500ms/step, ~10 steps out + ~10 back ≈ 10s walk; rest is grazing pause.
const (
	sheepMaxSteps  = 10
	sheepGrazePause = 21 * time.Second / gameSpeed
)

// SheepUnit wanders randomly across passable terrain near its home farm.
// It does not follow roads — sheep roam fields. On returning home it deposits
// one Wool into the farm's stockpile.
type SheepUnit struct {
	loc       cbb.Coord
	home      cbb.Coord
	stockpile *cbb.Inventory
	state     sheepState
	steps     int
	route     []cbb.Coord // return path home
	ts        time.Time
	prev      cbb.Coord // avoid immediate backtracking
	rng       *rand.Rand
}

func (s *SheepUnit) GetLoc() cbb.Coord          { return s.loc }
func (s *SheepUnit) Sprite() *ebiten.Image      { return sheepWorkerSprite }
func (s *SheepUnit) CanUpdate(t time.Time) bool { return t.After(s.ts) }

func (s *SheepUnit) Update(world cbb.World) {
	s.ts = time.Now().Add(500 * time.Millisecond)

	switch s.state {
	case sheepWandering:
		if s.steps >= sheepMaxSteps {
			s.startReturn(world)
			return
		}
		next, ok := s.nextTerrainTile(world)
		if !ok {
			s.startReturn(world)
			return
		}
		s.prev = s.loc
		s.loc = next
		s.steps++

	case sheepReturning:
		if len(s.route) <= 1 {
			// Arrived home — deposit wool and begin next graze after a pause.
			s.loc = s.home
			s.stockpile.Add(Wool, 1)
			s.steps = 0
			s.prev = s.home
			s.state = sheepWandering
			s.ts = time.Now().Add(sheepGrazePause)
			return
		}
		s.loc = s.route[len(s.route)-2]
		s.route = s.route[:len(s.route)-1]
	}
}

// startReturn finds a path home via the terrain tilemap and switches state.
// Falls back to teleport if pathfinding fails (e.g. sheep wandered into a dead end).
func (s *SheepUnit) startReturn(world cbb.World) {
	route, err := cbb.FindRoute(world.Tilemap(), s.loc, s.home)
	if err != nil || len(route) <= 1 {
		// Can't path home — teleport rather than get stuck.
		s.loc = s.home
		s.stockpile.Add(Wool, 1)
		s.steps = 0
		s.state = sheepWandering
		s.ts = time.Now().Add(sheepGrazePause)
		return
	}
	s.route = route
	s.state = sheepReturning
}

// nextTerrainTile picks a random passable terrain neighbour, preferring tiles
// that aren't the one we just came from (avoids boring back-and-forth).
func (s *SheepUnit) nextTerrainTile(world cbb.World) (cbb.Coord, bool) {
	tm := world.Tilemap()
	var forward, all []cbb.Coord
	for _, adj := range adjacentCoords(s.loc) {
		t, ok := tm.Tiles[adj]
		if !ok || !t.Passable {
			continue
		}
		all = append(all, adj)
		if adj != s.prev {
			forward = append(forward, adj)
		}
	}
	pool := forward
	if len(pool) == 0 {
		pool = all
	}
	if len(pool) == 0 {
		return cbb.Coord{}, false
	}
	return pool[s.rng.Intn(len(pool))], true
}
