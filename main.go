package main

import (
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	loadSprites()
	pixelgl.Run(run)
}
