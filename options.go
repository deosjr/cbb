package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type option struct {
	name         string
	buildingtype string // building or road
	radius       int    // 0 means no radius
	key          pixelgl.Button
	sprite       *pixel.Sprite
	newfunc      Instantiate
}

func getOptions() []option {
	return []option{
		{
			name:         "road",
			buildingtype: "road",
			key:          pixelgl.KeyR,
			sprite:       roadSprite,
		},
		{
			name:         "producer",
			buildingtype: "building",
			radius:       4,
			key:          pixelgl.KeyP,
			sprite:       producerSprite,
			newfunc:      NewProducer,
		},
		{
			name:         "consumer",
			buildingtype: "building",
			radius:       10,
			key:          pixelgl.KeyC,
			sprite:       consumerSprite,
			newfunc:      NewConsumer,
		},
	}
}
