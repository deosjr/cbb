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

// placedBuilding records everything needed to redraw a building each frame in
// isometric mode, where depth-sorted rendering requires per-frame drawing rather
// than a pre-baked batch image.
type placedBuilding struct {
	anchor    Coord
	sprite    *ebiten.Image
	footH     int // >0: combined multi-tile sprite, drawn via drawMultiBuildingToScreen
	w, h      int // effective footprint dimensions after rotation
}

// depth returns the isometric depth of the building's front (south-east) corner.
// Used for back-to-front draw ordering.
func (pb placedBuilding) depth() float64 {
	return pb.anchor.X + float64(pb.w-1) + pb.anchor.Y + float64(pb.h-1)
}

// Game is the main engine. Create it with NewGame and pass it to ebiten.RunGame.
type Game struct {
	world World

	camPos       Coord
	camSpeed     float64
	CamZoom      float64
	camZoomSpeed float64

	options   []Option
	selection Option
	rotation  int
	start     *Coord

	buildings  []placedBuilding // iso mode: sorted+redrawn each frame
	units      []Unit
	updatables []Updatable

	previewRoute  []Coord
	prevMouseTile Coord
	canPlace      bool

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
		canPlace:      true,
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
			if g.selection.Key != o.Key {
				g.rotation = 0
			}
			g.selection = o
			g.start = nil
			g.previewRoute = nil
			break
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) && g.selection.Kind == KindBuilding {
		g.rotation = (g.rotation + 1) % 4
	}

	mx, my := ebiten.CursorPosition()
	worldX, worldY := g.screenToWorld(float64(mx), float64(my))
	mouseTile := tileCoord(worldX, worldY)

	// Cache placement validity for cursor colouring in Draw.
	if g.selection.Kind == KindBuilding && g.selection.NewFunc != nil {
		sw, sh := g.effectiveSize()
		g.canPlace = g.footprintPassable(mouseTile, sw, sh)
		if g.canPlace {
			obj := g.selection.NewFunc()
			if p, ok := obj.(Placeable); ok {
				g.canPlace = p.CanPlace(mouseTile, g.world)
			}
		}
	} else {
		g.canPlace = true
	}

	// TODO: should not be able to build roads over houses!
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
							// Skip tiles occupied by a building footprint.
							if t, ok := g.world.Roads().Tiles[c]; ok && !t.Passable {
								continue
							}
							drawTileToBatch(tileBatch, c, roadSprite)
							g.world.Roads().Tiles[c] = Tile{Passable: true}
						}
					}
				}
			case KindBuilding:
				sw, sh := g.effectiveSize()
				if !g.footprintPassable(mouseTile, sw, sh) {
					break
				}
				obj := g.selection.NewFunc()
				if p, ok := obj.(Placeable); ok && !p.CanPlace(mouseTile, g.world) {
					break
				}
				units := obj.WhenPlaced(mouseTile, g.world)
				if r, ok := obj.(Rotatable); ok {
					r.SetRotation(g.rotation)
				}
				if isometricMode {
					pb := placedBuilding{anchor: mouseTile, w: sw, h: sh}
					if fps, ok := obj.(FootprintSpriteGetter); ok {
						pb.sprite, pb.footH = fps.GetFootprintSprite()
					} else {
						pb.sprite = obj.Sprite()
					}
					g.buildings = append(g.buildings, pb)
				} else {
					for dy := 0; dy < sh; dy++ {
						for dx := 0; dx < sw; dx++ {
							tc := Coord{mouseTile.X + float64(dx), mouseTile.Y + float64(dy)}
							drawTileToBatch(buildingBatch, tc, obj.Sprite())
						}
					}
				}
				for dy := 0; dy < sh; dy++ {
					for dx := 0; dx < sw; dx++ {
						tc := Coord{mouseTile.X + float64(dx), mouseTile.Y + float64(dy)}
						// Mark footprint as occupied but NOT passable: units cannot
						// path through building interiors, only via the access point.
						g.world.Roads().Tiles[tc] = Tile{Passable: false}
					}
				}
				ap := BuildingAccessPoint(mouseTile, sw, sh, g.rotation)
				if t, ok := g.world.Tilemap().Tiles[ap]; ok && t.Passable {
					// Access point overrides the footprint entry if it falls inside,
					// and is the sole road-connected entry/exit for the building.
					g.world.Roads().Tiles[ap] = Tile{Passable: true}
				}
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

	if isometricMode {
		// Depth-sort buildings back-to-front so foreground buildings correctly
		// occlude background ones. We sort a copy to preserve placement order.
		sorted := make([]placedBuilding, len(g.buildings))
		copy(sorted, g.buildings)
		slices.SortFunc(sorted, func(a, b placedBuilding) int {
			da, db := a.depth(), b.depth()
			if da < db {
				return -1
			}
			if da > db {
				return 1
			}
			return 0
		})
		for _, pb := range sorted {
			if pb.footH > 0 {
				drawMultiBuildingToScreen(screen, pb.anchor, pb.sprite, pb.footH, cam, nil)
			} else {
				for dy := 0; dy < pb.h; dy++ {
					for dx := 0; dx < pb.w; dx++ {
						tc := Coord{pb.anchor.X + float64(dx), pb.anchor.Y + float64(dy)}
						drawTileToScreen(screen, tc, pb.sprite, cam, nil)
					}
				}
			}
		}
	} else {
		screen.DrawImage(buildingBatch, camOp)
	}

	mx, my := ebiten.CursorPosition()
	worldX, worldY := g.screenToWorld(float64(mx), float64(my))
	mouseTile := tileCoord(worldX, worldY)

	for _, c := range g.previewRoute {
		drawTileToScreen(screen, c, roadSprite, cam, nil)
	}

	cursorSprite := g.selection.Sprite
	if g.selection.Kind == KindRoad {
		drawTileToScreen(screen, mouseTile, roadSprite, cam, nil)
	} else {
		var cs *ebiten.ColorScale
		if !g.canPlace {
			var red ebiten.ColorScale
			red.Scale(1, 0, 0, 1)
			cs = &red
		}
		sw, sh := g.effectiveSize()
		for dy := 0; dy < sh; dy++ {
			for dx := 0; dx < sw; dx++ {
				tc := Coord{mouseTile.X + float64(dx), mouseTile.Y + float64(dy)}
				drawTileToScreen(screen, tc, cursorSprite, cam, cs)
			}
		}
		if g.selection.NewFunc != nil {
			ap := BuildingAccessPoint(mouseTile, sw, sh, g.rotation)
			drawTileToScreen(screen, ap, highlightSprite, cam, nil)
		}
		if g.selection.Radius > 0 {
			cx := mouseTile.X + float64(sw-1)/2
			cy := mouseTile.Y + float64(sh-1)/2
			g.drawRadiusHighlight(screen, Coord{cx, cy}, cam)
		}
	}

	units := g.units
	if isometricMode {
		units = sortedByDepth(units)
	}
	for _, u := range units {
		drawTileToScreen(screen, u.GetLoc(), u.Sprite(), cam, nil)
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
			drawTileToScreen(screen, Coord{x, y}, highlightSprite, cam, nil)
		}
	}
}

// selectionSize returns the footprint of an option, treating 0 as 1.
func selectionSize(o Option) (w, h int) {
	w, h = o.SizeW, o.SizeH
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	return
}

// effectiveSize returns the footprint dimensions accounting for rotation.
// Odd rotations (west/east) swap width and height.
func (g *Game) effectiveSize() (w, h int) {
	w, h = selectionSize(g.selection)
	if g.rotation%2 == 1 {
		w, h = h, w
	}
	return
}

// footprintPassable reports whether all tiles in a w×h rectangle anchored at
// top-left corner are passable terrain AND not already occupied by a road or
// building.
func (g *Game) footprintPassable(anchor Coord, w, h int) bool {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			tc := Coord{anchor.X + float64(dx), anchor.Y + float64(dy)}
			if t, ok := g.world.Tilemap().Tiles[tc]; !ok || !t.Passable {
				return false
			}
			if _, occupied := g.world.Roads().Tiles[tc]; occupied {
				return false
			}
		}
	}
	return true
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

// AddUpdatable registers an Updatable that will be ticked each frame.
// Use this at game startup to add world-level simulation objects (e.g. a
// population ticker) that are not tied to a placed building.
func (g *Game) AddUpdatable(u Updatable) {
	g.updatables = append(g.updatables, u)
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
