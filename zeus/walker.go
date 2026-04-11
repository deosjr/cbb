package main

import (
	"math/rand"
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

type walkerState int

const (
	walkerWaiting    walkerState = iota // resting at home between patrols
	walkerPatrolling                    // wandering roads, recording coverage
	walkerReturning                     // following a route back to home
)

// ServiceWalker is a unit that wanders the road network from its home building,
// recording service coverage on every tile it visits. This is the key difference
// from Anno's cart-based delivery: houses don't receive physical goods; they
// become "served" whenever a walker passes nearby.
//
// Behaviour loop:
//  1. Wait at home for waitDur (cooldown between patrols).
//  2. Wander randomly along roads for up to maxSteps tiles.
//  3. Find a route home and walk back.
//  4. Repeat.
//
// Every tile the walker stands on is stamped in zeusWorld.coverage with the
// current time. Four cardinal neighbours are also stamped so that houses just
// off the road gain coverage too.
type ServiceWalker struct {
	loc         cbb.Coord
	home        cbb.Coord
	serviceType ServiceType
	sprite      *ebiten.Image

	state    walkerState
	steps    int
	maxSteps int
	route    []cbb.Coord // return path
	waitEnd  time.Time
	waitDur  time.Duration
	ts       time.Time

	rng  *rand.Rand
	prev cbb.Coord // last tile, to avoid immediate backtracking
}

// NewServiceWalker creates a ServiceWalker that starts its first patrol immediately.
func NewServiceWalker(home cbb.Coord, svc ServiceType, sprite *ebiten.Image, maxSteps int, waitDur time.Duration) *ServiceWalker {
	return &ServiceWalker{
		loc:         home,
		home:        home,
		serviceType: svc,
		sprite:      sprite,
		state:       walkerWaiting,
		maxSteps:    maxSteps,
		waitDur:     waitDur,
		waitEnd:     time.Now(), // start first patrol immediately
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (sw *ServiceWalker) GetLoc() cbb.Coord          { return sw.loc }
func (sw *ServiceWalker) Sprite() *ebiten.Image      { return sw.sprite }
func (sw *ServiceWalker) CanUpdate(t time.Time) bool { return t.After(sw.ts) }

func (sw *ServiceWalker) Update(world cbb.World) {
	sw.ts = time.Now().Add(400 * time.Millisecond)
	zw := world.(*zeusWorld)

	switch sw.state {
	case walkerWaiting:
		if time.Now().Before(sw.waitEnd) {
			return
		}
		sw.loc = sw.home
		sw.steps = 0
		sw.prev = sw.home
		sw.state = walkerPatrolling
		sw.stamp(zw)

	case walkerPatrolling:
		if sw.steps >= sw.maxSteps {
			sw.startReturn(world)
			return
		}
		next, ok := sw.nextRoadTile(world)
		if !ok {
			sw.startReturn(world)
			return
		}
		sw.prev = sw.loc
		sw.loc = next
		sw.steps++
		sw.stamp(zw)

	case walkerReturning:
		if len(sw.route) <= 1 {
			sw.loc = sw.home
			sw.waitEnd = time.Now().Add(sw.waitDur)
			sw.state = walkerWaiting
			return
		}
		sw.loc = sw.route[len(sw.route)-2]
		sw.route = sw.route[:len(sw.route)-1]
		sw.stamp(zw)
	}
}

// startReturn computes a route home and switches to walkerReturning.
// If pathfinding fails (e.g. walker is somehow off the road network), it teleports.
func (sw *ServiceWalker) startReturn(world cbb.World) {
	route, err := cbb.FindRoute(world.Roads(), sw.loc, sw.home)
	if err != nil || len(route) <= 1 {
		sw.loc = sw.home
		sw.waitEnd = time.Now().Add(sw.waitDur)
		sw.state = walkerWaiting
		return
	}
	sw.route = route
	sw.state = walkerReturning
}

// stamp records coverage at the current tile and its four neighbours.
// Neighbours get coverage so that houses adjacent to a road are served
// without needing to be directly on the road.
func (sw *ServiceWalker) stamp(zw *zeusWorld) {
	zw.RecordCoverage(sw.loc, sw.serviceType)
	for _, adj := range adjacentCoords(sw.loc) {
		zw.RecordCoverage(adj, sw.serviceType)
	}
}

// nextRoadTile picks a random road-adjacent tile, preferring forward movement
// (avoids the immediately previous tile unless it's the only option).
func (sw *ServiceWalker) nextRoadTile(world cbb.World) (cbb.Coord, bool) {
	roads := world.Roads()
	var forward, all []cbb.Coord
	for _, adj := range adjacentCoords(sw.loc) {
		if _, ok := roads.Tiles[adj]; !ok {
			continue
		}
		all = append(all, adj)
		if adj != sw.prev {
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
	return pool[sw.rng.Intn(len(pool))], true
}
