package main

// Good is the resource type for this game.
// Add new constants here as the game grows.
type Good int

const (
	NoGood Good = iota
)

func (g Good) Name() string {
	return [...]string{"nothing"}[g]
}
