package cbb

// Inventory is a simple counted store of resources.
// Use it in buildings and world structs to track goods.
type Inventory struct {
	stocks map[Resource]int
}

func NewInventory() *Inventory {
	return &Inventory{stocks: map[Resource]int{}}
}

// Add deposits n units of r.
func (inv *Inventory) Add(r Resource, n int) {
	inv.stocks[r] += n
}

// Take withdraws n units of r. Returns false and does nothing if insufficient stock.
func (inv *Inventory) Take(r Resource, n int) bool {
	if inv.stocks[r] < n {
		return false
	}
	inv.stocks[r] -= n
	return true
}

// Has reports whether at least n units of r are available.
func (inv *Inventory) Has(r Resource, n int) bool {
	return inv.stocks[r] >= n
}

// Count returns the current stock of r.
func (inv *Inventory) Count(r Resource) int {
	return inv.stocks[r]
}
