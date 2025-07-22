package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"slices"
)

type WorldID = ID

const UndefinedWorldID WorldID = -1

type World interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type WorldManager struct {
	activeWorld WorldID
	worlds map[WorldID]World
}

func NewWorldManager() *WorldManager {
	return &WorldManager{
		activeWorld:  UndefinedWorldID,
		worlds:        make(map[WorldID]World, 0),
	}
}

func (wm *WorldManager) Add(world World) {
	if _, exists := wm.worlds[world.ID()]; exists {
		return
	}

	wm.worlds[world.ID()] = world
}

func (wm *WorldManager) SetActive(worldID WorldID) {
	if _, exists := wm.worlds[worldID]; !exists {
		return
	}

	wm.activeWorld = worldID
}

func (wm *WorldManager) Remove(world World) {
	if _, exists := wm.worlds[world.ID()]; !exists {
		return
	}

	delete(wm.worlds, world.ID())
	if wm.activeWorld == world.ID() {
		wm.activeWorld = UndefinedWorldID
	}
}

func (wm *WorldManager) getActive() (World, bool) {
	if wm.activeWorld == UndefinedWorldID {
		return nil, false
	}

	world, exists := wm.worlds[wm.activeWorld]
	if !exists {
		return nil, false
	}

	return world, true
}

func (wm *WorldManager) Update() error {
	world, exists := wm.getActive()
	if !exists {
		return nil
	}

	if err := world.Update(); err != nil {
		return fmt.Errorf("error updating active world: %w", err)
	}
	
	return nil
}

func (wm *WorldManager) Draw(screen *ebiten.Image) {
	world, exists := wm.getActive()
	if !exists {
		return
	}

	world.Draw(screen)
}

