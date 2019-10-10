package main

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

func showFloor(tick int, r *sdl.Renderer, floor [][]*Sprite) {
	for y := 0; y <= MAX_Y; y++ {
		for x := 0; x <= MAX_X; x++ {
			sprite := floor[y][x]

			drawSpriteAt(tick, r, sprite, int32(x), int32(y), 0)
		}
	}
}

var floorTiles []*Sprite

func chooseRandomFloorSprite() *Sprite {
	if floorTiles == nil {
		floorTiles = spritesWithTag("floor")
	}
	n := rand.Intn(len(floorTiles))
	return floorTiles[n]
}

func generateFloor(width int) [][]*Sprite {
	floor := make([][]*Sprite, width)
	for y := 0; y < width; y++ {
		floor[y] = make([]*Sprite, width)
		for x := 0; x < width; x++ {
			floor[y][x] = chooseRandomFloorSprite()
		}
	}
	return floor
}
