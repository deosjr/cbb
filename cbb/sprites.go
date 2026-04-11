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
// Use this for flat terrain tiles in isometric games.
func NewIsoDiamondSprite(c color.Color) *ebiten.Image {
	base := image.NewRGBA(image.Rect(0, 0, isoTileW, isoTileH))
	r, g, b, a := c.RGBA()
	rc := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	for y := 0; y < isoTileH; y++ {
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

// NewIsoBoxSprite creates an isometric box sprite for one tile footprint.
// The sprite is isoTileW × (isoTileH + wallH). wallH controls building height.
// The left face uses wallColor; the right face is shaded 70% for depth.
// drawTileToBatch/drawTileToScreen automatically lift the sprite by wallH so
// the box floor aligns with the tile grid.
func NewIsoBoxSprite(wallColor, roofColor color.Color, wallH int) *ebiten.Image {
	sprH := isoTileH + wallH
	base := image.NewRGBA(image.Rect(0, 0, isoTileW, sprH))

	toRGBA := func(c color.Color) color.RGBA {
		r, g, b, a := c.RGBA()
		return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	}
	roofC := toRGBA(roofColor)
	wc := toRGBA(wallColor)
	rightWallC := color.RGBA{
		uint8(int(wc.R) * 7 / 10),
		uint8(int(wc.G) * 7 / 10),
		uint8(int(wc.B) * 7 / 10),
		wc.A,
	}

	// Draw walls column by column so each face is a proper isometric parallelogram.
	//
	// Left face: column x in [0, isoTileW/2).
	//   The lower-left edge of the roof diamond sits at
	//   y = isoTileH/2 + x*isoTileH/isoTileW for each column.
	//   The wall hangs wallH pixels straight down from that edge.
	for x := 0; x < isoTileW/2; x++ {
		yTop := isoTileH/2 + x*isoTileH/isoTileW
		for y := yTop; y < yTop+wallH; y++ {
			base.SetRGBA(x, y, wc)
		}
	}
	// Right face: column x in [isoTileW/2, isoTileW).
	//   The lower-right edge sits at y = isoTileH/2 + (isoTileW-x)*isoTileH/isoTileW.
	for x := isoTileW / 2; x < isoTileW; x++ {
		yTop := isoTileH/2 + (isoTileW-x)*isoTileH/isoTileW
		for y := yTop; y < yTop+wallH; y++ {
			base.SetRGBA(x, y, rightWallC)
		}
	}

	// Draw roof diamond on top, overwriting the upper wall area.
	for y := 0; y < isoTileH; y++ {
		t := y
		if isoTileH-1-y < t {
			t = isoTileH - 1 - y
		}
		halfW := (t + 1) * isoTileW / isoTileH
		left := isoTileW/2 - halfW
		right := isoTileW/2 + halfW
		for x := left; x < right; x++ {
			base.SetRGBA(x, y, roofC)
		}
	}

	return ebiten.NewImageFromImage(base)
}

// NewIsoBoxSpriteMulti creates a single combined isometric box sprite for a w×h
// tile footprint. The sprite is (w+h)*32 × (w+h)*16+wallH pixels.
//
// The anchor tile (0,0) top vertex is at sprite x = h*32 (i.e. shifted right by
// (h-1)*32 from the sprite's left edge). Use drawMultiBuildingToBatch to position
// it correctly; it subtracts that offset so the anchor aligns with its tile.
//
// Roof: a hexagonal diamond covering the full w×h footprint.
// Left face (camera-left): wallColor, right face: 70% darker.
// Both faces are proper iso parallelograms derived column by column.
func NewIsoBoxSpriteMulti(wallColor, roofColor color.Color, wallH, w, h int) *ebiten.Image {
	sprW := (w + h) * isoTileW / 2
	sprH := (w+h)*isoTileH/2 + wallH
	base := image.NewRGBA(image.Rect(0, 0, sprW, sprH))

	toRGBA := func(c color.Color) color.RGBA {
		r, g, b, a := c.RGBA()
		return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	}
	roofC := toRGBA(roofColor)
	wc := toRGBA(wallColor)
	rightWallC := color.RGBA{
		uint8(int(wc.R) * 7 / 10),
		uint8(int(wc.G) * 7 / 10),
		uint8(int(wc.B) * 7 / 10),
		wc.A,
	}

	// The combined w×h footprint diamond has these key x positions:
	//   topVertX  = h * isoTileW/2   — top vertex of the diamond (anchor tile top)
	//   midX      = w * isoTileW/2   — bottom vertex x (deepest point)
	//
	// Upper boundary at column x:
	//   x <= topVertX : yUpper = h*isoTileH/2 - x * isoTileH/isoTileW
	//   x >  topVertX : yUpper = (x-topVertX) * isoTileH/isoTileW
	//
	// Lower boundary at column x (where walls start):
	//   x <= midX     : yLower = h*isoTileH/2 + x * isoTileH/isoTileW
	//   x >  midX     : yLower = (w+h)*isoTileH/2 - (x-midX) * isoTileH/isoTileW
	topVertX := h * isoTileW / 2
	midX := w * isoTileW / 2

	for x := 0; x < sprW; x++ {
		// upper boundary
		var yUpper int
		if x <= topVertX {
			yUpper = h*isoTileH/2 - x*isoTileH/isoTileW
		} else {
			yUpper = (x - topVertX) * isoTileH / isoTileW
		}
		// lower boundary
		var yLower int
		if x <= midX {
			yLower = h*isoTileH/2 + x*isoTileH/isoTileW
		} else {
			yLower = (w+h)*isoTileH/2 - (x-midX)*isoTileH/isoTileW
		}
		// roof fill
		for y := yUpper; y < yLower; y++ {
			base.SetRGBA(x, y, roofC)
		}
		// wall below lower boundary
		if x < midX {
			for y := yLower; y < yLower+wallH; y++ {
				base.SetRGBA(x, y, wc)
			}
		} else {
			for y := yLower; y < yLower+wallH; y++ {
				base.SetRGBA(x, y, rightWallC)
			}
		}
	}

	return ebiten.NewImageFromImage(base)
}

// drawMultiBuildingToBatch draws a multi-tile building sprite onto a batch.
// The sprite was created by NewIsoBoxSpriteMulti with footprint (w, h).
// footH is the effective footprint height (h after rotation), used to shift
// the sprite left so the anchor tile's diamond aligns with coord c.
func drawMultiBuildingToBatch(target *ebiten.Image, c Coord, sprite *ebiten.Image, footH int) {
	op := &ebiten.DrawImageOptions{}
	sx, sy := isoToScreen(c)
	xOff := float64((footH - 1) * isoTileW / 2)
	// wallH is encoded in the sprite: sprH = (w+h)*isoTileH/2 + wallH,
	// and sprW = (w+h)*isoTileW/2, so wallH = sprH - sprW/2.
	wallH := float64(sprite.Bounds().Dy() - sprite.Bounds().Dx()/2)
	op.GeoM.Translate(sx+float64(isoOffsetX)-xOff, sy-wallH)
	target.DrawImage(sprite, op)
}

func drawMultiBuildingToScreen(screen *ebiten.Image, c Coord, sprite *ebiten.Image, footH int, cam ebiten.GeoM, cs *ebiten.ColorScale) {
	op := &ebiten.DrawImageOptions{}
	sx, sy := isoToScreen(c)
	xOff := float64((footH - 1) * isoTileW / 2)
	wallH := float64(sprite.Bounds().Dy() - sprite.Bounds().Dx()/2)
	op.GeoM.Translate(sx+float64(isoOffsetX)-xOff, sy-wallH)
	op.GeoM.Concat(cam)
	if cs != nil {
		op.ColorScale = *cs
	}
	screen.DrawImage(sprite, op)
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
		heightOffset := float64(sprite.Bounds().Dy() - isoTileH)
		op.GeoM.Translate(sx+float64(isoOffsetX), sy-heightOffset)
	} else {
		op.GeoM.Translate(c.X*resolution, c.Y*resolution)
	}
	target.DrawImage(sprite, op)
}

func drawTileToScreen(screen *ebiten.Image, c Coord, sprite *ebiten.Image, cam ebiten.GeoM, cs *ebiten.ColorScale) {
	op := &ebiten.DrawImageOptions{}
	if isometricMode {
		sx, sy := isoToScreen(c)
		heightOffset := float64(sprite.Bounds().Dy() - isoTileH)
		op.GeoM.Translate(sx+float64(isoOffsetX), sy-heightOffset)
	} else {
		op.GeoM.Translate(c.X*resolution, c.Y*resolution)
	}
	op.GeoM.Concat(cam)
	if cs != nil {
		op.ColorScale = *cs
	}
	screen.DrawImage(sprite, op)
}
