package main

// Good is the resource type for Zeus goods.
type Good int

const (
	Wheat    Good = iota // produced by Wheat Farm
	Meat                 // produced by Hunting Lodge
	OliveOil             // produced by Olive Farm (future)
	Wine                 // produced by Vineyard (future)
)

func (g Good) Name() string {
	return [...]string{"Wheat", "Meat", "Olive Oil", "Wine"}[g]
}

// AllGoods lists every Good constant.
var AllGoods = []Good{Wheat, Meat, OliveOil, Wine}
