package cbb

import (
	"math"

	"github.com/deosjr/Pathfinding/path"
	"github.com/hajimehoshi/ebiten/v2"
)

// Coord is a tile grid position. X and Y are always integer-valued in practice.
type Coord struct {
	X, Y float64
}

// Tile holds per-tile data.
type Tile struct {
	Passable bool
	// Sprite overrides the default passable/impassable tile sprite for this tile.
	// If nil, the library falls back to its built-in sprites.
	Sprite *ebiten.Image
}

// TileMap is the world grid, and also implements path.Graph for A* pathfinding.
type TileMap struct {
	Tiles map[Coord]Tile
}

func (tm *TileMap) Neighbours(n path.Node) []path.Node {
	p := n.(Coord)
	x, y := p.X, p.Y
	cardinal := []Coord{
		{x - 1, y},
		{x, y - 1},
		{x, y + 1},
		{x + 1, y},
	}
	neighbours := []path.Node{}
	for _, c := range cardinal {
		t, ok := tm.Tiles[c]
		if !ok || !t.Passable {
			continue
		}
		neighbours = append(neighbours, c)
	}
	return neighbours
}

func (tm *TileMap) G(n, neighbour path.Node) float64 {
	return 1
}

func (tm *TileMap) H(n, neighbour path.Node) float64 {
	p, q := n.(Coord), neighbour.(Coord)
	dx := q.X - p.X
	dy := q.Y - p.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// FindRoute finds a path between two tile coords on the given map.
// Returns the route as []Coord, hiding the path.Node type from callers.
func FindRoute(tm *TileMap, from, to Coord) ([]Coord, error) {
	nodes, err := path.FindRoute(tm, from, to)
	if err != nil {
		return nil, err
	}
	route := make([]Coord, len(nodes))
	for i, n := range nodes {
		route[i] = n.(Coord)
	}
	return route, nil
}

func tileCoord(wx, wy float64) Coord {
	return Coord{math.Floor(wx / resolution), math.Floor(wy / resolution)}
}
