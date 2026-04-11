package main

import "github.com/deosjr/tiles/cbb"

// Task represents a unit of work: deliver a resource to a building.
// This is a game-level concept; the cbb engine has no knowledge of it.
type Task interface {
	Destination() cbb.Building
	Resource() cbb.Resource
	Amount() int
}

// Receiver is optionally implemented by buildings that accept task deliveries.
type Receiver interface {
	Receive(Task, *taskWorld)
}

// taskWorld wraps cbb.BaseWorld and adds a simple FIFO task queue.
// Buildings post tasks via AddTask; gatherer units claim them via ClaimTask.
type taskWorld struct {
	*cbb.BaseWorld
	tasks []Task
}

func newTaskWorld(tilemap *cbb.TileMap) *taskWorld {
	return &taskWorld{BaseWorld: cbb.NewBaseWorld(tilemap)}
}

func (w *taskWorld) AddTask(t Task) {
	w.tasks = append(w.tasks, t)
}

func (w *taskWorld) ClaimTask() (Task, bool) {
	if len(w.tasks) == 0 {
		return nil, false
	}
	t := w.tasks[0]
	w.tasks = w.tasks[1:]
	return t, true
}
