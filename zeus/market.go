package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const marketPlaceCost = 75

// marketWalkerSteps is how many road tiles the market lady walks before turning back.
const marketWalkerSteps = 40

// marketWalkerCooldown is the wait time between market lady patrols.
const marketWalkerCooldown = 20 * time.Second / gameSpeed

// Market spawns a market-lady ServiceWalker that distributes food coverage.
// In Zeus/Caesar III, the market lady carries physical food from the market's
// storage; here we simplify: food is consumed centrally by PopulationTick, and
// the walker's presence is what grants ServiceFood coverage to nearby houses.
type Market struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int
}

func NewMarket() cbb.Building { return &Market{} }

func (m *Market) GetLoc() cbb.Coord                       { return m.loc }
func (m *Market) Sprite() *ebiten.Image                   { return marketSprite }
func (m *Market) AccessPoint() cbb.Coord                  { return m.accessPt }
func (m *Market) GetFootprintSprite() (*ebiten.Image, int) { return m.isoSprite, m.isoFootH }

func (m *Market) SetRotation(r int) {
	m.rotation = r
	m.isoSprite, m.isoFootH = cbb.NewIsoBoxSpriteMultiRotated(marketWall, marketRoof, marketWallH, 1, 1, r)
	m.accessPt = cbb.BuildingAccessPoint(m.loc, 1, 1, r)
}

func (m *Market) CanPlace(loc cbb.Coord, world cbb.World) bool {
	zw := world.(*zeusWorld)
	return zw.terrainAt(loc) != Water && zw.gold >= marketPlaceCost
}

func (m *Market) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	m.loc = loc
	m.SetRotation(0)
	zw := world.(*zeusWorld)
	zw.gold -= marketPlaceCost
	// Spawn the market lady walker from the access point.
	walker := NewServiceWalker(
		m.accessPt,
		ServiceFood,
		marketWalkerSprite,
		marketWalkerSteps,
		marketWalkerCooldown,
	)
	return []cbb.Unit{walker}
}
