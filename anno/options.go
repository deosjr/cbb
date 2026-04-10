package main

import (
	"github.com/deosjr/tiles/cbb"
	"github.com/hajimehoshi/ebiten/v2"
)

func getOptions() []cbb.Option {
	return []cbb.Option{
		{
			Name: "road",
			Kind: cbb.KindRoad,
			Key:  ebiten.KeyR,
		},
		{
			Name:    "warehouse",
			Kind:    cbb.KindBuilding,
			Key:     ebiten.KeyW,
			Sprite:  warehouseSprite,
			NewFunc: NewWarehouse,
		},
		{
			Name:    "fisherman",
			Kind:    cbb.KindBuilding,
			Key:     ebiten.KeyF,
			Sprite:  fishermanSprite,
			NewFunc: NewFishermanHut,
		},
		{
			Name:    "forester",
			Kind:    cbb.KindBuilding,
			Key:     ebiten.KeyT,
			Sprite:  foresterSprite,
			NewFunc: NewForesterHut,
		},
		{
			Name:    "hunter",
			Kind:    cbb.KindBuilding,
			Key:     ebiten.KeyU,
			Sprite:  hunterSprite,
			NewFunc: NewHunterLodge,
		},
		{
			Name:    "sheep farm",
			Kind:    cbb.KindBuilding,
			Key:     ebiten.KeyS,
			Sprite:  sheepFarmSprite,
			NewFunc: NewSheepFarm,
		},
		{
			Name:    "weaver",
			Kind:    cbb.KindBuilding,
			Key:     ebiten.KeyV,
			Sprite:  weaverSprite,
			NewFunc: NewWeavingHut,
		},
		{
			Name:    "house",
			Kind:    cbb.KindBuilding,
			Radius:  3,
			Key:     ebiten.KeyY,
			Sprite:  houseSprite,
			NewFunc: NewPioneerHouse,
		},
	}
}
