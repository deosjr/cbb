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
			Name:    "granary",
			Kind:    cbb.KindBuilding,
			SizeW:   2, SizeH: 2,
			Key:     ebiten.KeyG,
			Sprite:  granarySprite,
			NewFunc: NewGranary,
		},
		{
			Name:    "wheat farm",
			Kind:    cbb.KindBuilding,
			SizeW:   2, SizeH: 2,
			Key:     ebiten.KeyF,
			Sprite:  farmSprite,
			NewFunc: NewWheatFarm,
		},
		{
			Name:    "hunting lodge",
			Kind:    cbb.KindBuilding,
			SizeW:   1, SizeH: 1,
			Key:     ebiten.KeyH,
			Sprite:  hunterSprite,
			NewFunc: NewHuntingLodge,
		},
		{
			Name:    "market",
			Kind:    cbb.KindBuilding,
			SizeW:   1, SizeH: 1,
			Radius:  10, // visualise coverage radius in placement mode
			Key:     ebiten.KeyM,
			Sprite:  marketSprite,
			NewFunc: NewMarket,
		},
		{
			Name:    "fountain",
			Kind:    cbb.KindBuilding,
			SizeW:   1, SizeH: 1,
			Radius:  8,
			Key:     ebiten.KeyU,
			Sprite:  fountainSprite,
			NewFunc: NewFountain,
		},
		{
			Name:    "sanctuary",
			Kind:    cbb.KindBuilding,
			SizeW:   2, SizeH: 1,
			Radius:  12,
			Key:     ebiten.KeyS,
			Sprite:  sanctuarySprite,
			NewFunc: NewSanctuary,
		},
		{
			Name:    "house",
			Kind:    cbb.KindBuilding,
			SizeW:   1, SizeH: 1,
			Key:     ebiten.KeyY,
			Sprite:  houseSprites[0],
			NewFunc: NewGreekHouse,
		},
	}
}
