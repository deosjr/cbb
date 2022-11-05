package main

import (
	"math"

	"github.com/deosjr/Pathfinding/path"
	"github.com/faiface/pixel"
)

type coord pixel.Vec

type tile struct {
	passable bool
}

type TileMap struct {
	tiles map[coord]tile
}

func (tm *TileMap) Neighbours(n path.Node) []path.Node {
	p := n.(coord)
	x, y := p.X, p.Y
	cardinal := []coord{
		{x - 1, y},
		{x, y - 1},
		{x, y + 1},
		{x + 1, y},
	}
	points := []path.Node{}
	for _, c := range cardinal {
		t, ok := tm.tiles[c]
		if !ok || !t.passable {
			continue
		}
		points = append(points, c)
	}
	return points
}

func (tm *TileMap) G(n, neighbour path.Node) float64 {
	return 0
}

func (tm *TileMap) H(n, neighbour path.Node) float64 {
	p, q := n.(coord), neighbour.(coord)
	dx := float64(q.X - p.X)
	dy := float64(q.Y - p.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func tileCoord(v pixel.Vec) coord {
	res2 := resolution / 2
	dx := float64((int(v.X) + res2) / resolution)
	if int(v.X)+res2 < 0 {
		dx -= 1
	}
	dy := float64((int(v.Y) + res2) / resolution)
	if int(v.Y)+res2 < 0 {
		dy -= 1
	}
	return coord(pixel.V(dx, dy))
}

func middleVec(c coord) pixel.Vec {
	return pixel.V(float64(c.X)*resolution, float64(c.Y)*resolution)
}

func topLeftCornerVec(c coord) pixel.Vec {
	middle := middleVec(c)
	res2 := float64(resolution / 2)
	return pixel.V(middle.X-res2, middle.Y+res2)
}

func bottomRightCornerVec(c coord) pixel.Vec {
	middle := middleVec(c)
	res2 := float64(resolution / 2)
	return pixel.V(middle.X+res2, middle.Y-res2)
}
