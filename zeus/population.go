package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
)

const (
	// gameSpeed compresses real time: 10× means 1 in-game minute = 6 real seconds.
	gameSpeed = 10

	// cycleInterval is how often PopulationTick fires.
	cycleInterval = 60 * time.Second / gameSpeed
)

// PopulationTick fires every game cycle and handles house evolution, population
// growth/shrinkage, and tax collection for all houses at once.
//
// Zeus service model (unlike Anno inventory delivery):
//   - Food is consumed from the granary (1 unit per house per cycle).
//   - Services (hygiene, entertainment, etc.) are checked via the coverage map.
//   - If a house has all required services AND food, its population grows.
//   - If it's at the tier cap AND all services for the next tier are met, it evolves.
//   - If food OR a required service is missing, population shrinks and tiers devolve.
type PopulationTick struct {
	ts time.Time
}

func (pt *PopulationTick) CanUpdate(t time.Time) bool { return t.After(pt.ts) }

func (pt *PopulationTick) Update(world cbb.World) {
	pt.ts = time.Now().Add(cycleInterval)
	zw := world.(*zeusWorld)

	total := 0
	for _, h := range zw.houses {
		// Consume food: try wheat first, then meat.
		fed := zw.granary.Take(Wheat, 1) || zw.granary.Take(Meat, 1)

		// Check services for current tier.
		serviced := fed && allServicesCovered(zw, h.loc, h.tier)

		cap := tierCap[h.tier]

		if serviced {
			if h.pop < cap {
				h.pop++
			}
			// Attempt evolution when at cap and next tier's services are available.
			if h.pop >= cap && int(h.tier) < len(tierCap)-1 {
				nextTier := h.tier + 1
				if fed && allServicesCovered(zw, h.loc, nextTier) {
					h.tier = nextTier
					h.rebuildSprite()
				}
			}
		} else {
			if h.pop > 0 {
				h.pop--
			}
			// Devolve when empty and above base tier.
			if h.pop == 0 && h.tier > tierHovel {
				h.tier--
				h.rebuildSprite()
			}
		}

		total += h.pop
		zw.gold += h.pop // simplified tax: 1 gold per resident per cycle
	}
	zw.population = total
}

// allServicesCovered reports whether loc has fresh coverage for every service
// required at the given tier.
func allServicesCovered(zw *zeusWorld, loc cbb.Coord, tier houseTier) bool {
	for _, svc := range tierRequires[tier] {
		if !zw.HasCoverage(loc, svc) {
			return false
		}
	}
	return true
}
