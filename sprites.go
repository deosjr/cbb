package main

import (
	"image/color"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	producerSprite *ebiten.Image
	consumerSprite *ebiten.Image
	unitSprite     *ebiten.Image
)

func loadGameSprites() {
	producerSprite = cbb.NewSolidSprite(color.RGBA{0x8b, 0x00, 0x00, 0xff}) // darkred
	consumerSprite = cbb.NewSolidSprite(color.RGBA{0xff, 0xd7, 0x00, 0xff}) // gold
	unitSprite = cbb.NewSolidSprite(color.RGBA{0x00, 0x00, 0x8b, 0xff})     // darkblue
}
