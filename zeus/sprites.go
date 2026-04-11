package main

import (
	"image/color"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

// rgb packs three bytes into a fully-opaque RGBA.
func rgb(r, g, b uint8) color.RGBA { return color.RGBA{r, g, b, 0xff} }

// Building palette variables.
var (
	granaryWall, granaryRoof   color.Color
	farmWall, farmRoof         color.Color
	hunterWall, hunterRoof     color.Color
	marketWall, marketRoof     color.Color
	fountainWall, fountainRoof color.Color
	sanctuaryWall, sanctuaryRoof color.Color

	// Houses have a palette per tier (hovel → house → townhouse → manor).
	houseWalls [4]color.Color
	houseRoofs [4]color.Color
)

const (
	granaryWallH   = 20
	farmWallH      = 14
	hunterWallH    = 14
	marketWallH    = 16
	fountainWallH  = 22
	sanctuaryWallH = 32
	houseWallH     = 18
)

func initColors() {
	granaryWall = rgb(0xD2, 0xB4, 0x8C) // tan
	granaryRoof = rgb(0x8B, 0x45, 0x13) // saddlebrown

	farmWall = rgb(0xF5, 0xDE, 0xB3) // wheat
	farmRoof = rgb(0xC8, 0xA0, 0x50) // golden

	hunterWall = rgb(0x8B, 0x69, 0x14) // dark goldenrod
	hunterRoof = rgb(0x5C, 0x40, 0x10)

	marketWall = rgb(0xFF, 0xD7, 0x00) // gold
	marketRoof = rgb(0xCC, 0x88, 0x00)

	fountainWall = rgb(0xAD, 0xD8, 0xE6) // light blue
	fountainRoof = rgb(0x40, 0x80, 0xC0)

	sanctuaryWall = rgb(0xF0, 0xF0, 0xE8) // marble white
	sanctuaryRoof = rgb(0xCC, 0xCC, 0x44) // gold roof

	// Tier 0: humble hovel
	houseWalls[0] = rgb(0xD2, 0xB4, 0x8C)
	houseRoofs[0] = rgb(0x80, 0x80, 0x80)
	// Tier 1: simple house (food service)
	houseWalls[1] = rgb(0xFF, 0xFF, 0xE0)
	houseRoofs[1] = rgb(0xCC, 0x44, 0x22)
	// Tier 2: townhouse (food + hygiene)
	houseWalls[2] = rgb(0xFF, 0xFF, 0xFF)
	houseRoofs[2] = rgb(0x44, 0x66, 0xAA)
	// Tier 3: manor (food + hygiene + entertainment)
	houseWalls[3] = rgb(0xF8, 0xF8, 0xFF)
	houseRoofs[3] = rgb(0x22, 0x44, 0x88)
}

var (
	waterSprite  *ebiten.Image
	plainsSprite *ebiten.Image
	forestSprite *ebiten.Image

	granarySprite  *ebiten.Image
	farmSprite     *ebiten.Image
	hunterSprite   *ebiten.Image
	marketSprite   *ebiten.Image
	fountainSprite *ebiten.Image
	sanctuarySprite *ebiten.Image

	// houseSprites[tier] is the cursor sprite for each tier.
	houseSprites [4]*ebiten.Image

	// Walker unit sprites.
	marketWalkerSprite   *ebiten.Image
	fountainWalkerSprite *ebiten.Image
	sanctuaryWalkerSprite *ebiten.Image
	hunterWorkerSprite   *ebiten.Image
	farmWorkerSprite     *ebiten.Image
)

func loadSprites() {
	initColors()

	// Terrain
	waterSprite  = cbb.NewIsoDiamondSprite(color.RGBA{0x1E, 0x90, 0xFF, 0xff})
	plainsSprite = cbb.NewIsoDiamondSprite(color.RGBA{0xD2, 0xC8, 0x8C, 0xff}) // sandy/Mediterranean
	forestSprite = cbb.NewIsoDiamondSprite(color.RGBA{0x22, 0x8B, 0x22, 0xff})

	// Building cursor sprites (single tile, used in placement preview).
	granarySprite   = cbb.NewIsoBoxSprite(granaryWall, granaryRoof, granaryWallH)
	farmSprite      = cbb.NewIsoBoxSprite(farmWall, farmRoof, farmWallH)
	hunterSprite    = cbb.NewIsoBoxSprite(hunterWall, hunterRoof, hunterWallH)
	marketSprite    = cbb.NewIsoBoxSprite(marketWall, marketRoof, marketWallH)
	fountainSprite  = cbb.NewIsoBoxSprite(fountainWall, fountainRoof, fountainWallH)
	sanctuarySprite = cbb.NewIsoBoxSprite(sanctuaryWall, sanctuaryRoof, sanctuaryWallH)

	for i := range houseSprites {
		houseSprites[i] = cbb.NewIsoBoxSprite(houseWalls[i], houseRoofs[i], houseWallH)
	}

	// Walker unit sprites (small solid squares).
	marketWalkerSprite    = cbb.NewSolidSprite(color.RGBA{0xFF, 0xD7, 0x00, 0xff}) // gold
	fountainWalkerSprite  = cbb.NewSolidSprite(color.RGBA{0x40, 0x80, 0xC0, 0xff}) // blue
	sanctuaryWalkerSprite = cbb.NewSolidSprite(color.RGBA{0xFF, 0xFF, 0xFF, 0xff}) // white
	hunterWorkerSprite    = cbb.NewSolidSprite(color.RGBA{0x8B, 0x69, 0x14, 0xff})
	farmWorkerSprite      = cbb.NewSolidSprite(color.RGBA{0xC8, 0xA0, 0x50, 0xff})
}

func terrainSprite(t Terrain) *ebiten.Image {
	switch t {
	case Water:
		return waterSprite
	case Forest:
		return forestSprite
	default:
		return plainsSprite
	}
}
