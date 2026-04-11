package main

import (
	"math"
	"math/rand"

	"github.com/deosjr/tiles/cbb"
)

type Terrain int

const (
	Water  Terrain = iota
	Plains Terrain = iota
	Forest Terrain = iota
)

const (
	mapW = 80
	mapH = 80
)

func generateMap(seed int64) (*cbb.TileMap, map[cbb.Coord]Terrain) {
	rng := rand.New(rand.NewSource(seed))
	terrain := map[cbb.Coord]Terrain{}
	tiles := map[cbb.Coord]cbb.Tile{}

	cx, cy := float64(mapW)/2, float64(mapH)/2
	maxRadius := float64(mapW) / 2.5

	for y := 0; y < mapH; y++ {
		for x := 0; x < mapW; x++ {
			c := cbb.Coord{X: float64(x), Y: float64(y)}
			dx, dy := float64(x)-cx, float64(y)-cy
			dist := math.Sqrt(dx*dx + dy*dy)
			// Irregular coastline via per-tile noise
			threshold := maxRadius * (0.65 + 0.35*rng.Float64())

			if dist < threshold {
				t := Plains
				// Forest in the interior (~30% of land)
				if dist < maxRadius*0.55 && rng.Float64() < 0.30 {
					t = Forest
				}
				terrain[c] = t
				tiles[c] = cbb.Tile{Passable: true, Sprite: terrainSprite(t)}
			} else {
				terrain[c] = Water
				tiles[c] = cbb.Tile{Passable: false, Sprite: waterSprite}
			}
		}
	}

	return &cbb.TileMap{Tiles: tiles}, terrain
}

// adjacentCoords returns the 4 cardinal neighbours of c.
func adjacentCoords(c cbb.Coord) [4]cbb.Coord {
	return [4]cbb.Coord{
		{c.X - 1, c.Y},
		{c.X + 1, c.Y},
		{c.X, c.Y - 1},
		{c.X, c.Y + 1},
	}
}

// findNearestTerrain returns the closest tile of the given type from from.
func findNearestTerrain(zw *zeusWorld, from cbb.Coord, t Terrain) (cbb.Coord, bool) {
	best := cbb.Coord{}
	bestDist := math.MaxFloat64
	found := false
	for c, terrain := range zw.terrain {
		if terrain != t {
			continue
		}
		dx, dy := c.X-from.X, c.Y-from.Y
		if d := dx*dx + dy*dy; d < bestDist {
			bestDist = d
			best = c
			found = true
		}
	}
	return best, found
}
