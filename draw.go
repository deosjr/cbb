package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/deosjr/Pathfinding/path"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// TODO: proper injection
var (
	tasks   = []Task{}
	tilemap *TileMap
	roads   *TileMap
)

type option struct {
	name         string
	buildingtype string // building or road
	radius       int    // 0 means no radius
	key          pixelgl.Button
	sprite       *pixel.Sprite
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "TILES",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
	)

	frames := 0
	second := time.Tick(time.Second)

	var options = []option{
		{
			name:         "road",
			buildingtype: "road",
			radius:       0,
			key:          pixelgl.KeyR,
			sprite:       roadSprite,
		},
		{
			name:         "producer",
			buildingtype: "building",
			radius:       4,
			key:          pixelgl.KeyP,
			sprite:       producerSprite,
		},
		{
			name:         "consumer",
			buildingtype: "building",
			radius:       10,
			key:          pixelgl.KeyC,
			sprite:       consumerSprite,
		},
	}

	var (
		start, end *coord
		selection  = options[0]
	)

	// create the map and draw once; it is static
	m := map[coord]tile{}
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			m[coord(pixel.V(float64(x), float64(y)))] = tile{rand.Float64() < 0.6}
		}
	}
	tilemap = &TileMap{tiles: m}
	for c, t := range tilemap.tiles {
		if t.passable {
			drawTile(batch, c, passableSprite)
			continue
		}
		drawTile(batch, c, impassableSprite)
	}
	roads = &TileMap{tiles: map[coord]tile{}}

	updatables := []Updatable{}
	dynamicDraws := []Updatable{}

	camPos = pixel.V(400, 400)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		mousepos := cam.Unproject(win.MousePosition())
		mouseTilePos := tileCoord(mousepos)
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if t, ok := tilemap.tiles[mouseTilePos]; ok && t.passable {
				switch selection.name {
				case "road":
					switch {
					case start == nil:
						start = &mouseTilePos
					case end == nil:
						end = &mouseTilePos
						if start != nil && end != nil {
							route, err := path.FindRoute(tilemap, *start, *end)
							start = nil
							end = nil
							if err != nil {
								fmt.Println(err)
							}
							for _, n := range route {
								c := n.(coord)
								drawTile(batch, c, roadSprite)
								roads.tiles[c] = tile{passable: true}
							}
						}
					}
				case "producer":
					drawTile(batch2, mouseTilePos, producerSprite)
					obj := &producer{}
					obj.WhenPlaced(mouseTilePos)
					updatables = append(updatables, obj)
					roads.tiles[mouseTilePos] = tile{passable: true}
				case "consumer":
					drawTile(batch2, mouseTilePos, consumerSprite)
					obj := &consumer{}
					obj.WhenPlaced(mouseTilePos)
					updatables = append(updatables, obj.unit)
					dynamicDraws = append(dynamicDraws, obj.unit)
					roads.tiles[mouseTilePos] = tile{passable: true}
				}
			}
		}
		zoomScrollFactor := 1.0 / camZoom
		if win.Pressed(pixelgl.KeyH) || win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * zoomScrollFactor * dt
		}
		if win.Pressed(pixelgl.KeyL) || win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * zoomScrollFactor * dt
		}
		if win.Pressed(pixelgl.KeyJ) || win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * zoomScrollFactor * dt
		}
		if win.Pressed(pixelgl.KeyK) || win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * zoomScrollFactor * dt
		}
		if win.Pressed(pixelgl.KeyQ) {
			return
		}
		for _, o := range options {
			if win.Pressed(o.key) {
				selection = o
				start, end = nil, nil
				break
			}
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		for _, object := range updatables {
			if !object.CanUpdate(time.Now()) {
				continue
			}
			object.Update()
		}

		// actually drawing things to screen
		win.Clear(colornames.Black)
		batch.Draw(win)
		batch2.Draw(win)
		if mousepos != win.MousePreviousPosition() {
			batch3.Clear()
			highlightsprite := selection.sprite
			if selection.name == "road" && start != nil {
				route, err := path.FindRoute(tilemap, *start, mouseTilePos)
				if err == nil {
					for _, n := range route {
						c := n.(coord)
						drawTile(batch3, c, highlightsprite)
					}
				}
			}
			drawTile(batch3, mouseTilePos, highlightsprite)
			if selection.radius > 0 {
				vec := bottomRightCornerVec(mouseTilePos)
				radius := float64(selection.radius * resolution)
				c := pixel.C(vec, radius)
				//for y := vec.Y-radius+1; y < vec.Y+radius; y++{
				//for x := vec.X-radius+1; x < vec.X+radius; x++{
				for y := mouseTilePos.Y - radius; y <= mouseTilePos.Y+radius; y++ {
					for x := mouseTilePos.X - radius; x <= mouseTilePos.X+radius; x++ {
						check := coord(pixel.V(x, y))
						if !c.Contains(middleVec(check)) {
							continue
						}
						drawTile(batch3, check, highlightSprite)
					}
				}
			}
		}
		batch3.Draw(win)
		batch4.Clear()
		for _, obj := range dynamicDraws {
			drawTile(batch4, obj.(Unit).GetLoc(), unitSprite)
		}
		batch4.Draw(win)
		win.Update()

		// frame counter
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

// we print grid tiles at center but find coords by lefthandcorner
func drawTile(target pixel.Target, c coord, sprite *pixel.Sprite) {
	sprite.Draw(target, pixel.IM.Moved(middleVec(c)))
}
