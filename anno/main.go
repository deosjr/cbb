package main

import (
	"fmt"
	"time"

	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// gameSpeed compresses all production and consumption timers so the game
// is observable in real time. A value of 10 means 10 in-game minutes pass
// per real second, preserving the supply/demand ratios from the Anno data.
const gameSpeed = 10

// AnnoGame wraps cbb.Game to overlay the Anno-specific HUD.
type AnnoGame struct {
	*cbb.Game
	world *annoWorld
}

func (ag *AnnoGame) Draw(screen *ebiten.Image) {
	ag.Game.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"Gold: %d   Population: %d\n"+
			"Warehouse — Food:%d  Wood:%d  Wool:%d  Cloth:%d\n"+
			"R=road  W=warehouse  F=fisherman  T=forester  U=hunter  S=sheep  V=weaver  Y=house",
		ag.world.gold,
		ag.world.population,
		ag.world.warehouse.Count(Food),
		ag.world.warehouse.Count(Wood),
		ag.world.warehouse.Count(Wool),
		ag.world.warehouse.Count(Cloth),
	))
}

func main() {
	loadSprites()

	tilemap, terrain := generateMap(42)
	world := newAnnoWorld(tilemap, terrain)

	game := cbb.NewGame(world, getOptions(), true)
	game.CamZoom = 0.5
	game.AddUpdatable(&PopulationTick{ts: time.Now()})

	ebiten.SetWindowSize(cbb.ScreenW, cbb.ScreenH)
	ebiten.SetWindowTitle("Anno 1602")
	if err := ebiten.RunGame(&AnnoGame{Game: game, world: world}); err != nil {
		panic(err)
	}
}
