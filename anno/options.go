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
			Name:    "woodcutter",
			Kind:    cbb.KindBuilding,
			Key:     ebiten.KeyT,
			Sprite:  woodcutterSprite,
			NewFunc: NewWoodcutterHut,
		},
		{
			Name:    "fisherman",
			Kind:    cbb.KindBuilding,
			Key:     ebiten.KeyF,
			Sprite:  fishermanSprite,
			NewFunc: NewFishermanHut,
		},
		{
			Name:    "house",
			Kind:    cbb.KindBuilding,
			Radius:  3,
			Key:     ebiten.KeyY,
			Sprite:  houseSprite,
			NewFunc: NewSettlerHouse,
		},
	}
}
