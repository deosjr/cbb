package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
)

// cycleInterval is the real-time duration of one Anno game cycle (1 game minute).
// All production and consumption rates from the data are expressed per cycle.
const cycleInterval = 60 * time.Second / gameSpeed

// PopulationTick fires once per game cycle and handles all houses in one pass:
//   - aggregates total food consumption across every house
//   - deducts the total from the warehouse in a single Take call
//   - grows or shrinks each house's population based on whether food was available
//   - collects tax
//
// A fractional accumulator (foodDebt) means the correct per-capita rate from
// the Anno data (0.95 food / 80 people / cycle) is preserved at any population
// size without per-house timers.
type PopulationTick struct {
	ts       time.Time
	foodDebt float64 // sub-unit food consumption carried forward between cycles
}

func (pt *PopulationTick) CanUpdate(t time.Time) bool { return t.After(pt.ts) }

func (pt *PopulationTick) Update(world cbb.World) {
	pt.ts = time.Now().Add(cycleInterval)
	aw := world.(*annoWorld)

	// Sum population across all houses.
	totalPop := 0
	for _, h := range aw.houses {
		totalPop += h.pop
	}

	// Accumulate fractional food demand; deduct whole units when they build up.
	// Rate: 0.95 food per 80 people per cycle.
	pt.foodDebt += float64(totalPop) * 0.95 / 80.0
	foodNeeded := int(pt.foodDebt)
	pt.foodDebt -= float64(foodNeeded)

	// Single warehouse deduction for the whole population.
	fed := foodNeeded == 0 || aw.warehouse.Take(Food, foodNeeded)

	// Update each house and collect tax.
	for _, h := range aw.houses {
		if fed {
			if h.pop < maxPioneerPop {
				h.pop++
				aw.population++
			}
		} else if h.pop > 0 {
			h.pop--
			aw.population--
		}
		aw.gold += h.pop // 1 gold per resident per cycle (simplified)
	}
}
