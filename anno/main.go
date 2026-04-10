package main

import (
	"fmt"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// AnnoGame wraps cbb.Game to add the Anno-specific HUD.
type AnnoGame struct {
	*cbb.Game
	world *annoWorld
}

func (ag *AnnoGame) Draw(screen *ebiten.Image) {
	ag.Game.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"Gold: %d   Population: %d\nR=road  W=warehouse  T=woodcutter  F=fisherman  Y=house",
		ag.world.gold, ag.world.population,
	))
}

func main() {
	loadSprites()

	tilemap, terrain := generateMap(42)
	world := newAnnoWorld(tilemap, terrain)
	world.warehouse.Add(Wood, 10)
	world.warehouse.Add(Fish, 10)

	game := cbb.NewGame(world, getOptions())
	game.CamZoom = 0.5

	ebiten.SetWindowSize(cbb.ScreenW, cbb.ScreenH)
	ebiten.SetWindowTitle("Anno 1602")
	if err := ebiten.RunGame(&AnnoGame{Game: game, world: world}); err != nil {
		panic(err)
	}
}
