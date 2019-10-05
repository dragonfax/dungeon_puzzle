package main

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

func drawEntities(tick int, r *sdl.Renderer, entities []*PlacedEntity) {
	for _, entity := range entities {
		drawSpriteAt(tick, r, entity.Sprite, entity.Shape.X, entity.Shape.Y)
	}
}

type PlacedEntity struct {
	Sprite *Sprite
	Shape  *resolv.Rectangle
}

func removePlacedEntity(input []*PlacedEntity, entity *PlacedEntity) []*PlacedEntity {
	for i, e := range input {
		if e == entity {
			// delete the item
			return append(input[:i], input[i+1:]...)
		}
	}
	panic("can't remove entity")
}
