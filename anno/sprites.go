package main

import (
	"image/color"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	waterSprite   *ebiten.Image
	plainsSprite  *ebiten.Image
	forestSprite  *ebiten.Image

	warehouseSprite  *ebiten.Image
	woodcutterSprite *ebiten.Image
	fishermanSprite  *ebiten.Image
	houseSprite      *ebiten.Image

	woodWorkerSprite *ebiten.Image
	fishWorkerSprite *ebiten.Image
	carrierSprite    *ebiten.Image
)

func loadSprites() {
	waterSprite = cbb.NewSolidSprite(color.RGBA{0x1E, 0x90, 0xFF, 0xff})  // dodger blue
	plainsSprite = cbb.NewSolidSprite(color.RGBA{0x90, 0xEE, 0x90, 0xff}) // light green
	forestSprite = cbb.NewSolidSprite(color.RGBA{0x22, 0x8B, 0x22, 0xff}) // forest green

	warehouseSprite = cbb.NewSolidSprite(color.RGBA{0x8B, 0x45, 0x13, 0xff})  // saddlebrown
	woodcutterSprite = cbb.NewSolidSprite(color.RGBA{0xA0, 0x52, 0x2D, 0xff}) // sienna
	fishermanSprite = cbb.NewSolidSprite(color.RGBA{0x00, 0x80, 0x80, 0xff})  // teal
	houseSprite = cbb.NewSolidSprite(color.RGBA{0xFF, 0xD7, 0x00, 0xff})      // gold

	woodWorkerSprite = cbb.NewSolidSprite(color.RGBA{0x8B, 0x00, 0x00, 0xff}) // darkred
	fishWorkerSprite = cbb.NewSolidSprite(color.RGBA{0x00, 0x00, 0x8B, 0xff}) // darkblue
	carrierSprite = cbb.NewSolidSprite(color.RGBA{0xFF, 0xFF, 0x00, 0xff})    // yellow
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
