package main

type Good int

const (
	Wood Good = iota
	Fish
)

func (g Good) Name() string {
	return [...]string{"Wood", "Fish"}[g]
}
