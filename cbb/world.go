package cbb

type world struct {
	tilemap *TileMap
	roads   *TileMap
	tasks   []Task
}

func newWorld(tilemap *TileMap) *world {
	return &world{
		tilemap: tilemap,
		roads:   &TileMap{Tiles: map[Coord]Tile{}},
	}
}

func (w *world) AddTask(t Task) {
	w.tasks = append(w.tasks, t)
}

func (w *world) ClaimTask() (Task, bool) {
	if len(w.tasks) == 0 {
		return nil, false
	}
	t := w.tasks[0]
	w.tasks = w.tasks[1:]
	return t, true
}

func (w *world) Roads() *TileMap {
	return w.roads
}
