package cbb

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	resolution = 32
	ScreenW    = 1024
	ScreenH    = 768
)

var (
	tileBatch, buildingBatch *ebiten.Image

	passableSprite   *ebiten.Image
	impassableSprite *ebiten.Image
	roadSprite       *ebiten.Image
	highlightSprite  *ebiten.Image
)

func initSprites(mapW, mapH int) {
	passableSprite = NewSolidSprite(color.RGBA{0x00, 0x64, 0x00, 0xff})   // darkgreen
	impassableSprite = NewSolidSprite(color.RGBA{0xa0, 0x52, 0x2d, 0xff}) // sienna
	roadSprite = NewSolidSprite(color.RGBA{0xa9, 0xa9, 0xa9, 0xff})       // darkgrey
	highlightSprite = NewSolidSprite(color.RGBA{0x10, 0x10, 0x00, 0x10})

	tileBatch = ebiten.NewImage(mapW*resolution, mapH*resolution)
	buildingBatch = ebiten.NewImage(mapW*resolution, mapH*resolution)
}

// NewSolidSprite creates a resolution×resolution image filled with a single color.
// Games use this to create building and unit sprites.
func NewSolidSprite(c color.Color) *ebiten.Image {
	img := ebiten.NewImage(resolution, resolution)
	img.Fill(c)
	return img
}

func drawTileToBatch(target *ebiten.Image, c Coord, sprite *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.X*resolution, c.Y*resolution)
	target.DrawImage(sprite, op)
}

func drawTileToScreen(screen *ebiten.Image, c Coord, sprite *ebiten.Image, cam ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.X*resolution, c.Y*resolution)
	op.GeoM.Concat(cam)
	screen.DrawImage(sprite, op)
}
