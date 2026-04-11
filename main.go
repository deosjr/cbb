package main

import (
	"math/rand"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	loadGameSprites()

	m := map[cbb.Coord]cbb.Tile{}
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			m[cbb.Coord{X: float64(x), Y: float64(y)}] = cbb.Tile{Passable: rand.Float64() < 0.6}
		}
	}

	world := newTaskWorld(&cbb.TileMap{Tiles: m})
	game := cbb.NewGame(world, getOptions(), false)
	ebiten.SetWindowSize(cbb.ScreenW, cbb.ScreenH)
	ebiten.SetWindowTitle("TILES")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
