package cbb

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game is the main engine. Create it with NewGame and pass it to ebiten.RunGame.
type Game struct {
	world *world

	camPos       Coord
	camSpeed     float64
	camZoom      float64
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

// NewGame initialises the engine with a tile map and a set of build options.
// It must be called before ebiten.RunGame.
func NewGame(tilemap *TileMap, options []Option) *Game {
	mapW, mapH := mapDimensions(tilemap)
	initSprites(mapW, mapH)

	w := newWorld(tilemap)

	for c, t := range tilemap.Tiles {
		if t.Passable {
			drawTileToBatch(tileBatch, c, passableSprite)
		} else {
			drawTileToBatch(tileBatch, c, impassableSprite)
		}
	}

	return &Game{
		world:         w,
		camPos:        Coord{float64(mapW*resolution) / 2, float64(mapH*resolution) / 2},
		camSpeed:      500.0,
		camZoom:       1.0,
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

	zoomScrollFactor := 1.0 / g.camZoom
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
	g.camZoom *= math.Pow(g.camZoomSpeed, wheelY)

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
		route, err := FindRoute(g.world.tilemap, *g.start, mouseTile)
		if err == nil {
			g.previewRoute = route
		} else {
			g.previewRoute = nil
		}
		g.prevMouseTile = mouseTile
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if t, ok := g.world.tilemap.Tiles[mouseTile]; ok && t.Passable {
			switch g.selection.Kind {
			case KindRoad:
				if g.start == nil {
					start := mouseTile
					g.start = &start
					g.previewRoute = nil
				} else {
					route, err := FindRoute(g.world.tilemap, *g.start, mouseTile)
					g.start = nil
					g.previewRoute = nil
					if err == nil {
						for _, c := range route {
							drawTileToBatch(tileBatch, c, roadSprite)
							g.world.roads.Tiles[c] = Tile{Passable: true}
						}
					}
				}
			case KindBuilding:
				obj := g.selection.NewFunc()
				units := obj.WhenPlaced(mouseTile, g.world)
				drawTileToBatch(buildingBatch, mouseTile, obj.Sprite())
				g.world.roads.Tiles[mouseTile] = Tile{Passable: true}
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

	for _, u := range g.units {
		drawTileToScreen(screen, u.GetLoc(), u.Sprite(), cam)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenW, ScreenH
}

func (g *Game) cameraGeoM() ebiten.GeoM {
	var m ebiten.GeoM
	m.Translate(-g.camPos.X, -g.camPos.Y)
	m.Scale(g.camZoom, g.camZoom)
	m.Translate(float64(ScreenW)/2, float64(ScreenH)/2)
	return m
}

func (g *Game) screenToWorld(sx, sy float64) (float64, float64) {
	wx := (sx-float64(ScreenW)/2)/g.camZoom + g.camPos.X
	wy := (sy-float64(ScreenH)/2)/g.camZoom + g.camPos.Y
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
