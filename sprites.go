package main

import (
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

// sprites are 32x32 by default
const resolution = 32

var (
	sprites                       []*pixel.Sprite
	batch, batch2, batch3, batch4 *pixel.Batch

	roadSprite, passableSprite, impassableSprite, consumerSprite, producerSprite, unitSprite, highlightSprite *pixel.Sprite

	// adapted from colornames.Yellow but with opacity
	// TODO: opaque drawing doesnt work yet
	highlightcolor = color.RGBA{0x10, 0x10, 0x00, 0x10}
)

func loadSprites() {
	spritesheet := generateSpriteSheet()
	sprites = generateSprites(spritesheet)
	batch = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	batch2 = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	batch3 = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	batch4 = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	roadSprite = sprites[0]
	passableSprite = sprites[1]
	impassableSprite = sprites[2]
	consumerSprite = sprites[3]
	producerSprite = sprites[4]
	unitSprite = sprites[5]
	highlightSprite = sprites[6]
}

func generateSpriteSheet() pixel.Picture {
	spriteColours := []color.RGBA{colornames.Darkgrey, colornames.Darkgreen, colornames.Sienna, colornames.Gold, colornames.Darkred, colornames.Darkblue, highlightcolor}
	img := image.NewRGBA(image.Rect(0, 0, resolution*len(spriteColours), resolution))
	for y := 0; y < resolution; y++ {
		for x := 0; x < resolution*len(spriteColours); x++ {
			img.Set(x, y, spriteColours[x/resolution])
		}
	}
	return pixel.PictureDataFromImage(img)
}

func generateSprites(spritesheet pixel.Picture) []*pixel.Sprite {
	sprites := []*pixel.Sprite{}
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += resolution {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += resolution {
			sprites = append(sprites, pixel.NewSprite(spritesheet, pixel.R(x, y, x+resolution, y+resolution)))
		}
	}
	return sprites
}
