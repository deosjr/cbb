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
			Name:    "producer",
			Kind:    cbb.KindBuilding,
			Radius:  4,
			Key:     ebiten.KeyP,
			Sprite:  producerSprite,
			NewFunc: NewProducer,
		},
		{
			Name:    "consumer",
			Kind:    cbb.KindBuilding,
			Radius:  10,
			Key:     ebiten.KeyC,
			Sprite:  consumerSprite,
			NewFunc: NewConsumer,
		},
	}
}
