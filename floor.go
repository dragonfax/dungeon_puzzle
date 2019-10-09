package main

import (
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

func showFloor(tick int, r *sdl.Renderer, floor [][]*Sprite) {
	for y := 0; y < MAX_Y; y++ {
		for x := 0; x < MAX_X; x++ {
			sprite := floor[y][x]

			drawSpriteAt(tick, r, sprite, int32(x), int32(y), 0)
		}
	}
}

var floorTiles []*Sprite

func chooseRandomFloorSprite() *Sprite {
	if floorTiles == nil {
		floorTiles = spritesWithTag("floor")
		fmt.Printf("found %d floor tiles\n", len(floorTiles))
		fmt.Printf("%v\n%v\n", floorTiles[0], floorTiles[1])
	}
	n := rand.Intn(len(floorTiles))
	// fmt.Printf("choosing floor tile %d %s\n", n, floorTiles[n].Name)
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
