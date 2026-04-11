package main

import (
	"image/color"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

// rgb packs three bytes into a fully-opaque RGBA.
func rgb(r, g, b uint8) color.RGBA { return color.RGBA{r, g, b, 0xff} }

// Per-building palette and wall height. Used for both cursor (single-tile)
// and placed (multi-tile) sprites so they always match.
var (
	warehouseWall, warehouseRoof color.Color
	fishermanWall, fishermanRoof color.Color
	foresterWall, foresterRoof   color.Color
	hunterWall, hunterRoof       color.Color
	sheepFarmWall, sheepFarmRoof color.Color
	weaverWall, weaverRoof       color.Color
	houseWall, houseRoof         color.Color
	chapelWall, chapelRoof       color.Color
)

const (
	warehouseWallH = 22
	fishermanWallH = 14
	foresterWallH  = 18
	hunterWallH    = 14
	sheepFarmWallH = 18
	weaverWallH    = 18
	houseWallH     = 18
	chapelWallH    = 28
)

func initColors() {
	warehouseWall = rgb(0x8B, 0x45, 0x13)
	warehouseRoof = rgb(0x5A, 0x2D, 0x0A)
	fishermanWall = rgb(0x00, 0x80, 0x80)
	fishermanRoof = rgb(0x00, 0x50, 0x50)
	foresterWall  = rgb(0x5C, 0x40, 0x33)
	foresterRoof  = rgb(0x22, 0x8B, 0x22)
	hunterWall    = rgb(0x8B, 0x69, 0x14)
	hunterRoof    = rgb(0x5C, 0x40, 0x10)
	sheepFarmWall = rgb(0xF5, 0xF5, 0xDC)
	sheepFarmRoof = rgb(0xCC, 0x55, 0x33)
	weaverWall    = rgb(0xDA, 0x70, 0xD6)
	weaverRoof    = rgb(0x90, 0x20, 0x90)
	houseWall     = rgb(0xFF, 0xE0, 0x80)
	houseRoof     = rgb(0xCC, 0x44, 0x22)
	chapelWall    = rgb(0xE0, 0xE0, 0xE0)
	chapelRoof    = rgb(0x80, 0x80, 0x90)
}

var (
	waterSprite  *ebiten.Image
	plainsSprite *ebiten.Image
	forestSprite *ebiten.Image

	warehouseSprite *ebiten.Image
	fishermanSprite *ebiten.Image
	foresterSprite  *ebiten.Image
	hunterSprite    *ebiten.Image
	sheepFarmSprite *ebiten.Image
	weaverSprite    *ebiten.Image
	houseSprite     *ebiten.Image
	chapelSprite    *ebiten.Image

	fishWorkerSprite    *ebiten.Image
	foresterWorkerSprite *ebiten.Image
	hunterWorkerSprite  *ebiten.Image
	sheepWorkerSprite   *ebiten.Image
	carrierSprite       *ebiten.Image
	processorCartSprite *ebiten.Image
)

func loadSprites() {
	initColors()
	// Terrain
	waterSprite = cbb.NewIsoDiamondSprite(color.RGBA{0x1E, 0x90, 0xFF, 0xff})  // dodger blue
	plainsSprite = cbb.NewIsoDiamondSprite(color.RGBA{0x90, 0xEE, 0x90, 0xff}) // light green
	forestSprite = cbb.NewIsoDiamondSprite(color.RGBA{0x22, 0x8B, 0x22, 0xff}) // forest green

	// Cursor sprites (single tile, used in the build menu / placement preview).
	// Placed buildings use NewIsoBoxSpriteMulti via GetFootprintSprite().
	warehouseSprite = cbb.NewIsoBoxSprite(warehouseWall, warehouseRoof, warehouseWallH)
	fishermanSprite = cbb.NewIsoBoxSprite(fishermanWall, fishermanRoof, fishermanWallH)
	foresterSprite  = cbb.NewIsoBoxSprite(foresterWall, foresterRoof, foresterWallH)
	hunterSprite    = cbb.NewIsoBoxSprite(hunterWall, hunterRoof, hunterWallH)
	sheepFarmSprite = cbb.NewIsoBoxSprite(sheepFarmWall, sheepFarmRoof, sheepFarmWallH)
	weaverSprite    = cbb.NewIsoBoxSprite(weaverWall, weaverRoof, weaverWallH)
	houseSprite     = cbb.NewIsoBoxSprite(houseWall, houseRoof, houseWallH)
	chapelSprite    = cbb.NewIsoBoxSprite(chapelWall, chapelRoof, chapelWallH)

	// Worker units (small solid squares so they stand out on iso tiles)
	fishWorkerSprite     = cbb.NewSolidSprite(color.RGBA{0x00, 0x00, 0x8B, 0xff}) // darkblue
	foresterWorkerSprite = cbb.NewSolidSprite(color.RGBA{0x5C, 0x40, 0x33, 0xff}) // dark brown
	hunterWorkerSprite   = cbb.NewSolidSprite(color.RGBA{0x8B, 0x69, 0x14, 0xff}) // dark goldenrod
	sheepWorkerSprite    = cbb.NewSolidSprite(color.RGBA{0xAA, 0xAA, 0xAA, 0xff}) // light grey
	carrierSprite        = cbb.NewSolidSprite(color.RGBA{0xFF, 0xFF, 0x00, 0xff}) // yellow
	processorCartSprite  = cbb.NewSolidSprite(color.RGBA{0xDA, 0x70, 0xD6, 0xff}) // orchid
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
