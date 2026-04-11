package main

import (
	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const housePlaceCost = 25

// houseTier represents the evolution level of a GreekHouse.
// Each tier has a higher population cap and stricter service requirements.
//
//	Tier 0 — Hovel:      2 pop  — no services needed
//	Tier 1 — House:      5 pop  — food coverage
//	Tier 2 — Townhouse:  9 pop  — food + hygiene
//	Tier 3 — Manor:     14 pop  — food + hygiene + entertainment
type houseTier int

const (
	tierHovel     houseTier = 0
	tierHouse     houseTier = 1
	tierTownhouse houseTier = 2
	tierManor     houseTier = 3
)

// tierCap is the maximum population per house at each tier.
var tierCap = [4]int{2, 5, 9, 14}

// tierRequires lists the services a house needs to maintain (and evolve to) a tier.
// A house that loses a required service will eventually devolve.
var tierRequires = [4][]ServiceType{
	{},                                          // Hovel
	{ServiceFood},                               // House
	{ServiceFood, ServiceHygiene},               // Townhouse
	{ServiceFood, ServiceHygiene, ServiceEntertainment}, // Manor
}

// GreekHouse is the basic residential building in Zeus. Its tier evolves (or
// devolves) each population cycle based on service coverage in its neighbourhood.
type GreekHouse struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int

	tier houseTier
	pop  int
}

func NewGreekHouse() cbb.Building { return &GreekHouse{} }

func (h *GreekHouse) GetLoc() cbb.Coord { return h.loc }

// Sprite returns the cursor/single-tile sprite matching the current tier.
func (h *GreekHouse) Sprite() *ebiten.Image { return houseSprites[h.tier] }

func (h *GreekHouse) AccessPoint() cbb.Coord                    { return h.accessPt }
func (h *GreekHouse) GetFootprintSprite() (*ebiten.Image, int)  { return h.isoSprite, h.isoFootH }

func (h *GreekHouse) SetRotation(r int) {
	h.rotation = r
	h.rebuildSprite()
	h.accessPt = cbb.BuildingAccessPoint(h.loc, 1, 1, r)
}

func (h *GreekHouse) CanPlace(loc cbb.Coord, world cbb.World) bool {
	zw := world.(*zeusWorld)
	return zw.terrainAt(loc) != Water && zw.gold >= housePlaceCost
}

func (h *GreekHouse) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	h.loc = loc
	h.SetRotation(0)
	zw := world.(*zeusWorld)
	zw.gold -= housePlaceCost
	zw.houses = append(zw.houses, h)
	return nil
}

// rebuildSprite recreates the iso box sprite for the current tier and rotation.
// Called from SetRotation and also from PopulationTick when the tier changes.
func (h *GreekHouse) rebuildSprite() {
	h.isoSprite, h.isoFootH = isoBoxMulti(houseWalls[h.tier], houseRoofs[h.tier], houseWallH, 1, 1, h.rotation)
}
