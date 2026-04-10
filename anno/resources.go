package main

// Good is the resource type for Anno goods.
type Good int

const (
	Food Good = iota // produced by Fisher, Hunter
	Wood             // produced by Forester; used for construction
	Wool             // produced by Sheep Farm
	Cloth            // produced by Weaving Hut (2 Wool → 1 Cloth)
)

func (g Good) Name() string {
	return [...]string{"Food", "Wood", "Wool", "Cloth"}[g]
}

// AllGoods lists every Good constant, used by carts when sweeping stockpiles.
var AllGoods = []Good{Food, Wood, Wool, Cloth}
