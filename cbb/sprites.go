package cbb

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	resolution = 32
	ScreenW    = 1024
	ScreenH    = 768

	isoTileW = 64
	isoTileH = 32
)

var (
	tileBatch, buildingBatch *ebiten.Image

	passableSprite   *ebiten.Image
	impassableSprite *ebiten.Image
	roadSprite       *ebiten.Image
	highlightSprite  *ebiten.Image

	isometricMode bool
	isoOffsetX    int
)

func initSprites(mapW, mapH int, iso bool) {
	isometricMode = iso
	if iso {
		passableSprite = NewIsoDiamondSprite(color.RGBA{0x00, 0x64, 0x00, 0xff})
		impassableSprite = NewIsoDiamondSprite(color.RGBA{0xa0, 0x52, 0x2d, 0xff})
		roadSprite = NewIsoDiamondSprite(color.RGBA{0xa9, 0xa9, 0xa9, 0xff})
		highlightSprite = NewIsoDiamondSprite(color.RGBA{0x10, 0x10, 0x00, 0x10})

		isoOffsetX = (mapH - 1) * isoTileW / 2
		batchW := (mapW + mapH) * isoTileW / 2
		batchH := (mapW + mapH) * isoTileH / 2
		tileBatch = ebiten.NewImage(batchW, batchH)
		buildingBatch = ebiten.NewImage(batchW, batchH)
	} else {
		passableSprite = NewSolidSprite(color.RGBA{0x00, 0x64, 0x00, 0xff})
		impassableSprite = NewSolidSprite(color.RGBA{0xa0, 0x52, 0x2d, 0xff})
		roadSprite = NewSolidSprite(color.RGBA{0xa9, 0xa9, 0xa9, 0xff})
		highlightSprite = NewSolidSprite(color.RGBA{0x10, 0x10, 0x00, 0x10})

		tileBatch = ebiten.NewImage(mapW*resolution, mapH*resolution)
		buildingBatch = ebiten.NewImage(mapW*resolution, mapH*resolution)
	}
}

// NewSolidSprite creates a resolution×resolution image filled with a single color.
// Use this for flat (top-down) games.
func NewSolidSprite(c color.Color) *ebiten.Image {
	img := ebiten.NewImage(resolution, resolution)
	img.Fill(c)
	return img
}

// NewIsoDiamondSprite creates a 64×32 diamond-shaped sprite filled with a single color.
// Use this for isometric games.
func NewIsoDiamondSprite(c color.Color) *ebiten.Image {
	base := image.NewRGBA(image.Rect(0, 0, isoTileW, isoTileH))
	r, g, b, a := c.RGBA()
	rc := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	for y := 0; y < isoTileH; y++ {
		// Distance from nearest horizontal edge (0-indexed)
		t := y
		if isoTileH-1-y < t {
			t = isoTileH - 1 - y
		}
		halfW := (t + 1) * isoTileW / isoTileH
		left := isoTileW/2 - halfW
		right := isoTileW/2 + halfW
		for x := left; x < right; x++ {
			base.SetRGBA(x, y, rc)
		}
	}
	return ebiten.NewImageFromImage(base)
}

// isoToScreen converts a tile grid coord to pixel position within the iso batch.
// The returned position is the top-left corner of the diamond, before isoOffsetX.
func isoToScreen(c Coord) (float64, float64) {
	sx := (c.X - c.Y) * float64(isoTileW/2)
	sy := (c.X + c.Y) * float64(isoTileH/2)
	return sx, sy
}

func drawTileToBatch(target *ebiten.Image, c Coord, sprite *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	if isometricMode {
		sx, sy := isoToScreen(c)
		op.GeoM.Translate(sx+float64(isoOffsetX), sy)
	} else {
		op.GeoM.Translate(c.X*resolution, c.Y*resolution)
	}
	target.DrawImage(sprite, op)
}

func drawTileToScreen(screen *ebiten.Image, c Coord, sprite *ebiten.Image, cam ebiten.GeoM, cs *ebiten.ColorScale) {
	op := &ebiten.DrawImageOptions{}
	if isometricMode {
		sx, sy := isoToScreen(c)
		op.GeoM.Translate(sx+float64(isoOffsetX), sy)
	} else {
		op.GeoM.Translate(c.X*resolution, c.Y*resolution)
	}
	op.GeoM.Concat(cam)
	if cs != nil {
		op.ColorScale = *cs
	}
	screen.DrawImage(sprite, op)
}
