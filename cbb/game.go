package cbb

import (
	"fmt"
	"image/color"
	"math"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game is the main engine. Create it with NewGame and pass it to ebiten.RunGame.
type Game struct {
	world World

	camPos       Coord
	camSpeed     float64
	CamZoom      float64
	camZoomSpeed float64

	options   []Option
	selection Option
	start     *Coord

	units      []Unit
	updatables []Updatable

	previewRoute  []Coord
	prevMouseTile Coord

	last          time.Time
	frames        int
	lastFPSUpdate time.Time
}

// NewGame initialises the engine with a world and a set of build options.
// The world must embed BaseWorld or otherwise implement the World interface.
// Call this before ebiten.RunGame.
func NewGame(world World, options []Option, isometric bool) *Game {
	tilemap := world.Tilemap()
	mapW, mapH := mapDimensions(tilemap)
	initSprites(mapW, mapH, isometric)

	for c, t := range tilemap.Tiles {
		sprite := t.Sprite
		if sprite == nil {
			if t.Passable {
				sprite = passableSprite
			} else {
				sprite = impassableSprite
			}
		}
		drawTileToBatch(tileBatch, c, sprite)
	}

	return &Game{
		world:         world,
		camPos:        Coord{float64(mapW*resolution) / 2, float64(mapH*resolution) / 2},
		camSpeed:      500.0,
		CamZoom:       1.0,
		camZoomSpeed:  1.2,
		options:       options,
		selection:     options[0],
		last:          time.Now(),
		lastFPSUpdate: time.Now(),
	}
}

func (g *Game) Update() error {
	dt := time.Since(g.last).Seconds()
	g.last = time.Now()

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	zoomScrollFactor := 1.0 / g.CamZoom
	if ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camPos.X -= g.camSpeed * zoomScrollFactor * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camPos.X += g.camSpeed * zoomScrollFactor * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camPos.Y += g.camSpeed * zoomScrollFactor * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camPos.Y -= g.camSpeed * zoomScrollFactor * dt
	}

	_, wheelY := ebiten.Wheel()
	g.CamZoom *= math.Pow(g.camZoomSpeed, wheelY)

	for _, o := range g.options {
		if ebiten.IsKeyPressed(o.Key) {
			g.selection = o
			g.start = nil
			g.previewRoute = nil
			break
		}
	}

	mx, my := ebiten.CursorPosition()
	worldX, worldY := g.screenToWorld(float64(mx), float64(my))
	mouseTile := tileCoord(worldX, worldY)

	if g.selection.Kind == KindRoad && g.start != nil && mouseTile != g.prevMouseTile {
		route, err := FindRoute(g.world.Tilemap(), *g.start, mouseTile)
		if err == nil {
			g.previewRoute = route
		} else {
			g.previewRoute = nil
		}
		g.prevMouseTile = mouseTile
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if t, ok := g.world.Tilemap().Tiles[mouseTile]; ok && t.Passable {
			switch g.selection.Kind {
			case KindRoad:
				if g.start == nil {
					start := mouseTile
					g.start = &start
					g.previewRoute = nil
				} else {
					route, err := FindRoute(g.world.Tilemap(), *g.start, mouseTile)
					g.start = nil
					g.previewRoute = nil
					if err == nil {
						for _, c := range route {
							drawTileToBatch(tileBatch, c, roadSprite)
							g.world.Roads().Tiles[c] = Tile{Passable: true}
						}
					}
				}
			case KindBuilding:
				obj := g.selection.NewFunc()
				if p, ok := obj.(Placeable); ok && !p.CanPlace(mouseTile, g.world) {
					break
				}
				units := obj.WhenPlaced(mouseTile, g.world)
				drawTileToBatch(buildingBatch, mouseTile, obj.Sprite())
				g.world.Roads().Tiles[mouseTile] = Tile{Passable: true}
				if u, ok := obj.(Updatable); ok {
					g.updatables = append(g.updatables, u)
				}
				for _, u := range units {
					g.units = append(g.units, u)
					g.updatables = append(g.updatables, u)
				}
			}
		}
	}

	now := time.Now()
	for _, obj := range g.updatables {
		if !obj.CanUpdate(now) {
			continue
		}
		obj.Update(g.world)
	}

	g.frames++
	if time.Since(g.lastFPSUpdate) >= time.Second {
		ebiten.SetWindowTitle(fmt.Sprintf("TILES | FPS: %d", g.frames))
		g.frames = 0
		g.lastFPSUpdate = time.Now()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	cam := g.cameraGeoM()
	camOp := &ebiten.DrawImageOptions{}
	camOp.GeoM = cam

	screen.DrawImage(tileBatch, camOp)
	screen.DrawImage(buildingBatch, camOp)

	mx, my := ebiten.CursorPosition()
	worldX, worldY := g.screenToWorld(float64(mx), float64(my))
	mouseTile := tileCoord(worldX, worldY)

	for _, c := range g.previewRoute {
		drawTileToScreen(screen, c, roadSprite, cam)
	}

	cursorSprite := g.selection.Sprite
	if g.selection.Kind == KindRoad {
		cursorSprite = roadSprite
	}
	drawTileToScreen(screen, mouseTile, cursorSprite, cam)

	if g.selection.Radius > 0 {
		g.drawRadiusHighlight(screen, mouseTile, cam)
	}

	units := g.units
	if isometricMode {
		units = sortedByDepth(units)
	}
	for _, u := range units {
		drawTileToScreen(screen, u.GetLoc(), u.Sprite(), cam)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenW, ScreenH
}

func (g *Game) cameraGeoM() ebiten.GeoM {
	var m ebiten.GeoM
	m.Translate(-g.camPos.X, -g.camPos.Y)
	m.Scale(g.CamZoom, g.CamZoom)
	m.Translate(float64(ScreenW)/2, float64(ScreenH)/2)
	return m
}

func (g *Game) screenToWorld(sx, sy float64) (float64, float64) {
	wx := (sx-float64(ScreenW)/2)/g.CamZoom + g.camPos.X
	wy := (sy-float64(ScreenH)/2)/g.CamZoom + g.camPos.Y
	return wx, wy
}

func (g *Game) drawRadiusHighlight(screen *ebiten.Image, mouseTile Coord, cam ebiten.GeoM) {
	cx, cy := mouseTile.X, mouseTile.Y
	radius := float64(g.selection.Radius)
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			dx, dy := x-cx, y-cy
			if math.Sqrt(dx*dx+dy*dy) > radius {
				continue
			}
			drawTileToScreen(screen, Coord{x, y}, highlightSprite, cam)
		}
	}
}

func sortedByDepth(units []Unit) []Unit {
	sorted := make([]Unit, len(units))
	copy(sorted, units)
	slices.SortFunc(sorted, func(a, b Unit) int {
		da := a.GetLoc().X + a.GetLoc().Y
		db := b.GetLoc().X + b.GetLoc().Y
		if da < db {
			return -1
		}
		if da > db {
			return 1
		}
		return 0
	})
	return sorted
}

func mapDimensions(tm *TileMap) (w, h int) {
	for c := range tm.Tiles {
		if int(c.X)+1 > w {
			w = int(c.X) + 1
		}
		if int(c.Y)+1 > h {
			h = int(c.Y) + 1
		}
	}
	return
}
