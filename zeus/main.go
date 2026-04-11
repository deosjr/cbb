package main

import (
	"fmt"
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ZeusGame wraps cbb.Game to overlay the Zeus-specific HUD.
type ZeusGame struct {
	*cbb.Game
	world *zeusWorld
}

func (zg *ZeusGame) Draw(screen *ebiten.Image) {
	zg.Game.Draw(screen)

	// Count covered houses per service for the HUD.
	foodCovered, hygieneCovered := 0, 0
	for _, h := range zg.world.houses {
		if zg.world.HasCoverage(h.loc, ServiceFood) {
			foodCovered++
		}
		if zg.world.HasCoverage(h.loc, ServiceHygiene) {
			hygieneCovered++
		}
	}
	total := len(zg.world.houses)

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"Gold: %d   Population: %d   Houses: %d\n"+
			"Granary — Wheat: %d  Meat: %d\n"+
			"Coverage — Food: %d/%d  Hygiene: %d/%d\n"+
			"R=road  G=granary  F=farm  H=hunter  M=market  U=fountain  S=sanctuary  Y=house\n"+
			"Tab=rotate",
		zg.world.gold,
		zg.world.population,
		total,
		zg.world.granary.Count(Wheat),
		zg.world.granary.Count(Meat),
		foodCovered, total,
		hygieneCovered, total,
	))
}

func main() {
	loadSprites()

	tilemap, terrain := generateMap(7)
	world := newZeusWorld(tilemap, terrain)

	game := cbb.NewGame(world, getOptions(), true)
	game.CamZoom = 0.5
	game.AddUpdatable(&PopulationTick{ts: time.Now()})

	ebiten.SetWindowSize(cbb.ScreenW, cbb.ScreenH)
	ebiten.SetWindowTitle("Zeus — City Builder")
	if err := ebiten.RunGame(&ZeusGame{Game: game, world: world}); err != nil {
		panic(err)
	}
}
