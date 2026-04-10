package main

import (
	"image/color"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

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

	fishWorkerSprite    *ebiten.Image
	foresterWorkerSprite *ebiten.Image
	hunterWorkerSprite  *ebiten.Image
	sheepWorkerSprite   *ebiten.Image
	carrierSprite       *ebiten.Image
	processorCartSprite *ebiten.Image
)

func loadSprites() {
	// Terrain
	waterSprite = cbb.NewIsoDiamondSprite(color.RGBA{0x1E, 0x90, 0xFF, 0xff})  // dodger blue
	plainsSprite = cbb.NewIsoDiamondSprite(color.RGBA{0x90, 0xEE, 0x90, 0xff}) // light green
	forestSprite = cbb.NewIsoDiamondSprite(color.RGBA{0x22, 0x8B, 0x22, 0xff}) // forest green

	// Buildings (iso diamond sprites)
	warehouseSprite = cbb.NewIsoDiamondSprite(color.RGBA{0x8B, 0x45, 0x13, 0xff}) // saddlebrown
	fishermanSprite = cbb.NewIsoDiamondSprite(color.RGBA{0x00, 0x80, 0x80, 0xff}) // teal
	foresterSprite  = cbb.NewIsoDiamondSprite(color.RGBA{0x5C, 0x40, 0x33, 0xff}) // dark brown
	hunterSprite    = cbb.NewIsoDiamondSprite(color.RGBA{0x8B, 0x69, 0x14, 0xff}) // dark goldenrod
	sheepFarmSprite = cbb.NewIsoDiamondSprite(color.RGBA{0xF5, 0xF5, 0xDC, 0xff}) // beige
	weaverSprite    = cbb.NewIsoDiamondSprite(color.RGBA{0xDA, 0x70, 0xD6, 0xff}) // orchid
	houseSprite     = cbb.NewIsoDiamondSprite(color.RGBA{0xFF, 0xD7, 0x00, 0xff}) // gold

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
