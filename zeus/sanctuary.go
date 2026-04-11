package main

import (
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

const sanctuaryPlaceCost = 200

// sanctuaryWalkerSteps is how many tiles the priest walker covers per procession.
const sanctuaryWalkerSteps = 50

// sanctuaryWalkerCooldown is the rest period between priestly processions.
const sanctuaryWalkerCooldown = 30 * time.Second / gameSpeed

// Sanctuary (temple) spawns a priest/procession walker providing ServiceReligion
// coverage. In the full Zeus game, specific gods require specific sanctuaries;
// here we use a single generic sanctuary to demonstrate the pattern.
//
// Houses need entertainment + religion to become Manors (tier 3).
// The entertainment service (ServiceEntertainment) would be provided by a
// Theater building — not yet implemented — so tier 3 evolution is future work.
type Sanctuary struct {
	loc       cbb.Coord
	rotation  int
	accessPt  cbb.Coord
	isoSprite *ebiten.Image
	isoFootH  int
}

func NewSanctuary() cbb.Building { return &Sanctuary{} }

func (s *Sanctuary) GetLoc() cbb.Coord                       { return s.loc }
func (s *Sanctuary) Sprite() *ebiten.Image                   { return sanctuarySprite }
func (s *Sanctuary) AccessPoint() cbb.Coord                  { return s.accessPt }
func (s *Sanctuary) GetFootprintSprite() (*ebiten.Image, int) { return s.isoSprite, s.isoFootH }

func (s *Sanctuary) SetRotation(r int) {
	s.rotation = r
	s.isoSprite, s.isoFootH = isoBoxMulti(sanctuaryWall, sanctuaryRoof, sanctuaryWallH, 2, 1, r)
	sw, sh := 2, 1
	if r%2 == 1 {
		sw, sh = 1, 2
	}
	s.accessPt = cbb.BuildingAccessPoint(s.loc, sw, sh, r)
}

func (s *Sanctuary) CanPlace(loc cbb.Coord, world cbb.World) bool {
	zw := world.(*zeusWorld)
	return zw.terrainAt(loc) != Water && zw.gold >= sanctuaryPlaceCost
}

func (s *Sanctuary) WhenPlaced(loc cbb.Coord, world cbb.World) []cbb.Unit {
	s.loc = loc
	s.SetRotation(0)
	zw := world.(*zeusWorld)
	zw.gold -= sanctuaryPlaceCost
	walker := NewServiceWalker(
		s.accessPt,
		ServiceReligion,
		sanctuaryWalkerSprite,
		sanctuaryWalkerSteps,
		sanctuaryWalkerCooldown,
	)
	return []cbb.Unit{walker}
}
